package config

// Region for the amazon services
const Region string = "REGION"

// GogetaSQSEndpoint is the aws sqs endpoint for the Gogeta queue
const GogetaSQSEndpoint string = "GOGETA_SQS_ENDPOINT"

// MrrobotSQSEndpoint is the aws sqs endpoint for the MrRobot queue
const MrrobotSQSEndpoint string = "MRROBOT_SQS_ENDPOINT"

// HalGogetaAPI is the api endpoint for hal monitoring gogeta
const HalGogetaAPI string = "HAL_GOGETA_API"

// HalMrRobotAPI is the api endpoint for hal monitoring MrRobot
const HalMrRobotAPI string = "HAL_MRROBOT_API"

// Auth0ClientSecret is the authentication token to use for microservices
const Auth0ClientSecret string = "AUTH0_CLIENT_SECRET"

// LogEndpoint is the endpoint for sending log data
const LogEndpoint string = "PAPERTRAIL_ENDPOINT"

// DevMode defines whether or not to run Mr.Robot in developer mode
const DevMode string = "DEV_MODE"
