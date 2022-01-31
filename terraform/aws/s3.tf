resource "aws_s3_bucket" "s3" {
  bucket = var.storage_name
  acl    = "private"

  tags = {
    Name = var.storage_name
  }
}