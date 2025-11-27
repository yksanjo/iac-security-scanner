# Example Terraform file with security issues for testing

resource "aws_s3_bucket" "public_bucket" {
  bucket = "my-public-bucket"
  acl    = "public-read"  # SECURITY ISSUE: Public access
}

resource "aws_db_instance" "public_db" {
  identifier     = "example-db"
  engine         = "mysql"
  instance_class = "db.t3.micro"
  publicly_accessible = true  # SECURITY ISSUE: Publicly accessible database
}

resource "aws_instance" "unencrypted" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"
  
  # SECURITY ISSUE: No encryption specified
  root_block_device {
    volume_type = "gp3"
    volume_size = 20
  }
}

resource "aws_security_group" "open_ssh" {
  name = "open-ssh"
  
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # SECURITY ISSUE: Open to internet
  }
}

