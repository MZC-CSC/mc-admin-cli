# 공통 변수 설정

variable "project_id" {
  description = "GCP 프로젝트 ID"
  type        = string
  default     = "spatial-conduit-399006"
}

variable "region" {
  description = "GCP 리전"
  type        = string
  default     = "asia-east1"
}
