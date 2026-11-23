# AWS CSP role mapping — bootstrap change and cleanup for existing installs

Applies to mc-iam-manager 0.7.3 and the matching mc-admin-cli bootstrap (`conf/docker/conf/mc-iam-manager/1_setup_auto.sh`, `init_aws_csp_config`).

## What changed

| Step | Before | Now |
|---|---|---|
| 1. CSP account (`aws-default`) | created | unchanged |
| 2. IdP Config (`aws-oidc`: `role_arn`, `oidc_provider_arn`, `audience`) | created | unchanged |
| 3. CspRole | registered the management role `${PREFIX}-platformadmin` as a CspRole | **removed** |
| 4. Mapping | mapped admin, operator, viewer, billadmin, billviewer to the management role | **replaced**: each Master Role `X` is mapped to `mciam-X` |

- The management role (the IdP Config `role_arn`) is the role IAM assumes to write to AWS. It must not be used for users. Every user who received temporary credentials through the old mapping got the management role's permissions (for example `IAMFullAccess`).
- mc-iam-manager 0.7.3 refuses to register, update, delete, or map a management role (HTTP 409). Removing a mapping (`DELETE /api/roles/csp-roles`) is still allowed, so existing mappings can be cleaned up.
- The new Step 3 calls `PUT /api/roles/id/{id}` for each role in `MC_IAM_MANAGER_PREDEFINED_ROLE`. It adds the `csp` sub and keeps the existing subs and description. It sends `cspRoles=[{cspRoleName: X, cspType: aws, authMethod: OIDC, cspIdpConfigId}]`.
  - IAM adds the `mciam-` prefix. `MC_IAM_MANAGER_CSP_ROLE_PREFIX` only affects the management role name, so the prefix is never doubled.
  - If `mciam-X` exists in AWS, IAM registers it and maps it. Otherwise IAM creates it with a trust policy built from the IdP Config. IAM checks the OIDC provider and audience before creating.
  - A failed item is reported as a warning and the script continues. Running the script again reports `unchanged`.

> **New AWS roles on a clean install.** A clean install now creates five roles in the AWS account: `mciam-admin`, `mciam-operator`, `mciam-viewer`, `mciam-billadmin` and `mciam-billviewer`. They have a trust policy only and **no permission policies**. Temporary credentials issued for them cannot do anything until permission policies are attached (CSP Roles screen, attachPolicy).

## Cleanup for existing installs (dev, stage)

Do these steps after upgrading mc-iam-manager to 0.7.3 and mc-admin-cli to this version.

### 1. Find the management role
```sql
SELECT id, name, config->>'role_arn' AS role_arn, is_active FROM mcmp_csp_idp_configs;
```
The role name is the last path element of `role_arn`, for example `mciam-platformadmin`. Some older environments use `mciam-admin`.

### 2. Remove mappings to the management role
- **Automatic**: running the new bootstrap (Step 3) issues a PUT for each predefined role. If a role still maps to the management role under the same name it requests (for example `admin` → `mciam-admin` when the management role is `mciam-admin`), IAM removes that mapping and reports `failed (protected, existing mapping removed)`. Other mappings to the management role are removed because the PUT replaces the CSP mappings of that role with `mciam-X`.
- **Manual** (roles that are not in `MC_IAM_MANAGER_PREDEFINED_ROLE`, or without re-running the bootstrap):
  ```sql
  -- list remaining mappings to the management role
  SELECT m.role_id, r.name AS master_role, m.csp_role_id, m.auth_method
    FROM mcmp_role_csp_role_mappings m
    JOIN mcmp_role_masters r ON r.id = m.role_id
    JOIN mcmp_role_csp_roles c ON c.id = m.csp_role_id
   WHERE c.iam_identifier = '<role_arn>' OR c.name = '<management role name>';
  ```
  Remove each one with the API (platformAdmin):
  ```bash
  curl -X DELETE "$IAM/api/roles/csp-roles" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
    -d '{"roleId":"<role_id>","cspRoleId":"<csp_role_id>","authMethod":"OIDC"}'
  ```

### 3. The management role's CspRole record
The old Step 3 left a CspRole record for the management role. Either keep it (IAM no longer lets anyone map, update, or delete it) or delete it with SQL after step 2:
```sql
DELETE FROM mcmp_role_csp_roles
 WHERE (iam_identifier = '<role_arn>' OR name = '<management role name>')
   AND NOT EXISTS (SELECT 1 FROM mcmp_role_csp_role_mappings m WHERE m.csp_role_id = mcmp_role_csp_roles.id);
```
**Do not delete it with the API** (`DELETE /api/roles/csp/id/{id}`). That API also deletes the AWS role. 0.7.3 refuses it with 409, but older versions would delete the management role and stop all IAM writes to AWS.

### 4. Check the management role's permission for the OIDC provider check
IAM 0.7.3 calls `iam:GetOpenIDConnectProvider` before creating a role. Without the permission it logs a WARN and skips the check.
```bash
aws iam simulate-principal-policy \
  --policy-source-arn arn:aws:iam::<account>:role/<management role> \
  --action-names iam:GetOpenIDConnectProvider \
  --resource-arns arn:aws:iam::<account>:oidc-provider/<domain>/auth/realms/<realm> \
  --query 'EvaluationResults[0].EvalDecision'
```
`allowed` means no change is needed. Roles with `IAMFullAccess` already allow it. Otherwise add this statement to the management role (approval required, it is a shared role):
```json
{"Effect": "Allow", "Action": "iam:GetOpenIDConnectProvider", "Resource": "arn:aws:iam::<account>:oidc-provider/*"}
```

### 5. Check the result
- Users get temporary credentials as `assumed-role/mciam-<role>/...`, not as the management role.
- Attach permission policies to the new `mciam-*` roles as needed.
