from __future__ import print_function
from boto3.session import Session
import boto3
import json
import os.path
import botocore
import tempfile
import zipfile
import traceback

from jinja2 import Environment, FileSystemLoader

code_pipeline = boto3.client('codepipeline')

class Template:
    def __init__(self, name):
        self.name = name

    def render(self, file, attributes):
        env = Environment(loader=FileSystemLoader('beanstalk'))
        template = env.get_template(self.name)
        return template.render(attributes)


def lambda_handler(event, context):
    """The Lambda function handler

    Args:
        event: The event passed by Lambda
        context: The context passed by Lambda

    """

    try:
        # Extract the Job ID
        job_id = event['CodePipeline.job']['id']

        # Extract the Job Data
        job_data = event['CodePipeline.job']['data']

        # Get the list of artifacts passed to the function
        artifacts = job_data['inputArtifacts']

        # Get the artifact details
        artifact_data = find_artifact(artifacts, 'MyAppBuild')

        # Get S3 client to access artifact with
        s3 = setup_s3_client(job_data)

        version_filename = 'build.version'
        version = get_file_from_artifact(s3, artifact_data, version_filename).strip()
        attributes = {
            'version': version
        }
        
        artifact_data = find_artifact(artifacts, 'MyInfra')
        template_file = 'Dockerrun.aws.json.j2'
        file_data = get_file_from_artifact(s3, artifact_data, template_file)

        rendered_data = Environment().from_string(file_data).render(attributes)
        output_artifact = generate_output_artifact({ 'Dockerrun.aws.json': rendered_data })

        outputArtifactDetails = find_artifact(job_data['outputArtifacts'], 'GenBuildParams')
        upload_artifact(s3, output_artifact, outputArtifactDetails)

        put_job_success(job_id, 'Configuration file created successfully')

    except Exception as e:
        # If any other exceptions which we didn't expect are raised
        # then fail the job and log the exception message.
        print('Function failed due to exception.')
        print(e)
        traceback.print_exc()
        put_job_failure(job_id, 'Function exception: ' + str(e))

def put_job_success(job, message):
    """Notify CodePipeline of a successful job

    Args:
        job: The CodePipeline job ID
        message: A message to be logged relating to the job status

    Raises:
        Exception: Any exception thrown by .put_job_success_result()

    """
    print('Putting job success')
    print(message)
    code_pipeline.put_job_success_result(jobId=job)

def upload_artifact(s3, file, artifact_location):
    bucket = artifact_location['location']['s3Location']['bucketName']
    key = artifact_location['location']['s3Location']['objectKey']
    with open(file.name, 'rb') as data:
        s3.upload_fileobj(data, bucket, key)

def generate_output_artifact(output_files):
    tmp_file = tempfile.NamedTemporaryFile(delete=False)
    zf = zipfile.ZipFile(tmp_file.name, mode='w')
    try:
        for f, data in output_files.iteritems():
            zf.writestr(f, data)
    finally:
        zf.close()

    return tmp_file

def find_artifact(artifacts, name):
    """Finds the artifact 'name' among the 'artifacts'

    Args:
        artifacts: The list of artifacts available to the function
        name: The artifact we wish to use
    Returns:
        The artifact dictionary found
    Raises:
        Exception: If no matching artifact is found

    """
    for artifact in artifacts:
        if artifact['name'] == name:
            return artifact

    raise Exception('Input artifact named "{0}" not found in event'.format(name))

def get_file_from_artifact(s3, artifact, file_in_zip):
    """Gets a file from a zipped artifact in S3

    Downloads the artifact from the S3 artifact store to a temporary file
    then extracts the zip and returns the file containing the file_in_zip

    Args:
        artifact: The artifact to download indicated as a combination of bucket and key
        file_in_zip: The path to the file within the zip containing the template

    Returns:
        The location of a temporary file on the filesystem?

    Raises:
        Exception: Any exception thrown while downloading the artifact or unzipping it

    """

    tmp_file = tempfile.NamedTemporaryFile()
    bucket = artifact['location']['s3Location']['bucketName']
    key = artifact['location']['s3Location']['objectKey']
    
    print("Downloading file {} from bucket {}".format(key, bucket))

    with tempfile.NamedTemporaryFile() as tmp_file:
        s3.download_file(bucket, key, tmp_file.name)
        with zipfile.ZipFile(tmp_file.name, 'r') as zip:
            return zip.read(file_in_zip)

def setup_s3_client(job_data):
    """Creates an S3 client

    Uses the credentials passed in the event by CodePipeline. These
    credentials can be used to access the artifact bucket.

    Args:
        job_data: The job data structure

    Returns:
        An S3 client with the appropriate credentials

    """
    key_id = job_data['artifactCredentials']['accessKeyId']
    key_secret = job_data['artifactCredentials']['secretAccessKey']
    session_token = job_data['artifactCredentials']['sessionToken']

    session = Session(aws_access_key_id=key_id,
        aws_secret_access_key=key_secret,
        aws_session_token=session_token)
    return session.client('s3', config=botocore.client.Config(signature_version='s3v4'))

if __name__ == '__main__':
    lambda_handler(json.load(open("tests/fixtures/render_event.json")), None)
