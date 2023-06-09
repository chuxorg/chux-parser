on:
  release:
    types:
      - created
jobs:
  terraform:
    name: 'Terraform'
    runs-on: ubuntu-latest
    env:
      AWS_REGION: ${{ secrets.AWS_REGION }}

    steps:
    - name: Checkout repository
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: '1.19'

    - name: Build Docker image
      run: |
        docker build -t chux-parser:latest .
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v1
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: ${{ secrets.AWS_REGION }} 
    
    - name: Login to Amazon ECR
      id: login-ecr
      uses: aws-actions/amazon-ecr-login@v1
      
    - name: Tag Docker image
      run: |
        docker tag chux-parser:latest ${{ steps.login-ecr.outputs.registry }}/chux-parser:${GITHUB_REF##*/}

    - name: Push Docker image to Amazon ECR
      run: |
        docker push ${{ steps.login-ecr.outputs.registry }}/chux-parser:${GITHUB_REF##*/}
        
    - name: Set up Terraform
      uses: hashicorp/setup-terraform@v1
      with:
        terraform_version: 1.0.11

    - name: Terraform Initialize
      run: terraform init
      working-directory: tf/
      env:
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        AWS_SDK_LOAD_CONFIG: 1
        TF_VAR_s3_bucket_name: chux-terraform-state
        TF_VAR_s3_key: chux-parser-terraform.tfstate
        TF_VAR_s3_region: ${{ secrets.AWS_REGION }}
        TF_VAR_s3_encrypt: true

    - name: Terraform Validate
      run: terraform validate
      working-directory: tf/
      env:
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        AWS_REGION: ${{ secrets.AWS_REGION }}
        AWS_SDK_LOAD_CONFIG: 1

    - name: Update Terraform variable with image URI
      run: |
        TAG_NAME=${GITHUB_REF##*/}
        echo "image_uri = \"${{ steps.login-ecr.outputs.registry }}/chux-parser:${TAG_NAME}\"" > tf/variables.auto.tfvars

    - name: Terraform Plan
      run: terraform plan -var="aws_access_key_id=${{ secrets.AWS_ACCESS_KEY_ID }}" -var="aws_secret_access_key=${{ secrets.AWS_SECRET_ACCESS_KEY }}" -var="aws_region=${{ secrets.AWS_REGION }}"
      working-directory: tf/
      env:
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        AWS_REGION: ${{ secrets.AWS_REGION }}
        AWS_SDK_LOAD_CONFIG: 1
    
    - name: Terraform Apply
      run: terraform apply -auto-approve -var="aws_access_key_id=${{ secrets.AWS_ACCESS_KEY_ID }}" -var="aws_secret_access_key=${{ secrets.AWS_SECRET_ACCESS_KEY }}" -var="aws_region=${{ secrets.AWS_REGION }}"
      working-directory: tf/
      env:
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        AWS_REGION: ${{ secrets.AWS_REGION }}
        AWS_SDK_LOAD_CONFIG: 1
