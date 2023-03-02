SCRIPT_DIRECTORY=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

ROOT_PROJECT_DIRECTORY="${SCRIPT_DIRECTORY}/.."
GENERATED_DECORATORS_EXTENTION=".gen.go"
GO_DECORATORS_TEMPLATE_DIRECTORY="${ROOT_PROJECT_DIRECTORY}/core/pkg/decorators/templates"
GO_CZAP_LOGGER_DECORATOR_TEMPLATE="${GO_DECORATORS_TEMPLATE_DIRECTORY}/czap_logger.template.go"
GO_OPENTELEMETRY_TRACING_DECORATOR_TEMPLATE="${GO_DECORATORS_TEMPLATE_DIRECTORY}/opentelemetry_tracing.template.go"
GO_PROMETHEUS_METRICS_DECORATOR_TEMPLATE="${GO_DECORATORS_TEMPLATE_DIRECTORY}/prometheus_metrics.template.go"

cd $ROOT_PROJECT_DIRECTORY

export GOPATH="/home/am/go"
export PATH="$PATH:$(go env GOPATH)/bin"

GO_WRAP_BINARY="${GOPATH}/bin/gowrap"

# user service
# user service: repositories
USER_SERVICE_REPOSITORY_INTERFACE_FOLDER="${ROOT_PROJECT_DIRECTORY}/user_service/repositories"
${GO_WRAP_BINARY} gen \
    -p ${USER_SERVICE_REPOSITORY_INTERFACE_FOLDER} \
    -i UserRepository \
    -t ${GO_CZAP_LOGGER_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/repositories/repositorydecorators/czap_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${USER_SERVICE_REPOSITORY_INTERFACE_FOLDER} \
    -i UserRepository \
    -t ${GO_OPENTELEMETRY_TRACING_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/repositories/repositorydecorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -v metricPrefix=database_query \
    -p ${USER_SERVICE_REPOSITORY_INTERFACE_FOLDER} \
    -i UserRepository \
    -t ${GO_PROMETHEUS_METRICS_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/repositories/repositorydecorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# user service: services
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/services \
    -i UserService \
    -t ${GO_CZAP_LOGGER_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/services/servicedecorators/czap_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/services \
    -i UserService \
    -t ${GO_OPENTELEMETRY_TRACING_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/services/servicedecorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -v metricPrefix=service_request \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/services \
    -i UserService \
    -t ${GO_PROMETHEUS_METRICS_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/services/servicedecorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# user service: controllers
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/controllers \
    -i UserController \
    -t ${GO_CZAP_LOGGER_DECORATOR_TEMPLATE} \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/controllers/controllerdecorators/czap_logger${GENERATED_DECORATORS_EXTENTION}


# sandbox service
# sandbox service: repositories
SANDBOX_REPOSITORY_INTERFACE_FOLDER="${ROOT_PROJECT_DIRECTORY}/sandbox_service/apientity"
SANDBOX_REPOSITORY_GENERATED_DECORATORS_FOLDER="${ROOT_PROJECT_DIRECTORY}/sandbox_service/decorators/repositorydecorators"
${GO_WRAP_BINARY} gen \
    -p ${SANDBOX_REPOSITORY_INTERFACE_FOLDER} \
    -i SandboxRepository \
    -t ${GO_CZAP_LOGGER_DECORATOR_TEMPLATE} \
    -o ${SANDBOX_REPOSITORY_GENERATED_DECORATORS_FOLDER}/czap_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${SANDBOX_REPOSITORY_INTERFACE_FOLDER} \
    -i SandboxRepository \
    -t ${GO_OPENTELEMETRY_TRACING_DECORATOR_TEMPLATE} \
    -o ${SANDBOX_REPOSITORY_GENERATED_DECORATORS_FOLDER}/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

# ${GO_WRAP_BINARY} gen \
#     -v metricPrefix=database_query \
#     -p ${SANDBOX_REPOSITORY_INTERFACE_FOLDER} \
#     -i SandboxRepository \
#     -t ${GO_PROMETHEUS_METRICS_DECORATOR_TEMPLATE} \
#     -o ${SANDBOX_REPOSITORY_GENERATED_DECORATORS_FOLDER}/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# sandbox service: service
SANDBOX_SERVICE_INTERFACE_FOLDER="${ROOT_PROJECT_DIRECTORY}/sandbox_service/apientity"
SANDBOX_SERVICE_GENERATED_DECORATORS_FOLDER="${ROOT_PROJECT_DIRECTORY}/sandbox_service/decorators/servicedecorators"
${GO_WRAP_BINARY} gen \
    -p ${SANDBOX_SERVICE_INTERFACE_FOLDER} \
    -i SandboxService \
    -t ${GO_CZAP_LOGGER_DECORATOR_TEMPLATE} \
    -o ${SANDBOX_SERVICE_GENERATED_DECORATORS_FOLDER}/czap_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${SANDBOX_SERVICE_INTERFACE_FOLDER} \
    -i SandboxService \
    -t ${GO_OPENTELEMETRY_TRACING_DECORATOR_TEMPLATE} \
    -o ${SANDBOX_SERVICE_GENERATED_DECORATORS_FOLDER}/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

# ${GO_WRAP_BINARY} gen \
#     -v metricPrefix=database_query \
#     -p ${SANDBOX_SERVICE_INTERFACE_FOLDER} \
#     -i SandboxService \
#     -t ${GO_PROMETHEUS_METRICS_DECORATOR_TEMPLATE} \
#     -o ${SANDBOX_SERVICE_GENERATED_DECORATORS_FOLDER}/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}


# sandbox service: service
SANDBOX_CONTROLLER_INTERFACE_FOLDER="${ROOT_PROJECT_DIRECTORY}/sandbox_service/apientity"
SANDBOX_CONTROLLER_GENERATED_DECORATORS_FOLDER="${ROOT_PROJECT_DIRECTORY}/sandbox_service/decorators/controllerdecorators"
${GO_WRAP_BINARY} gen \
    -p ${SANDBOX_CONTROLLER_INTERFACE_FOLDER} \
    -i SandboxController \
    -t ${GO_CZAP_LOGGER_DECORATOR_TEMPLATE} \
    -o ${SANDBOX_CONTROLLER_GENERATED_DECORATORS_FOLDER}/czap_logger${GENERATED_DECORATORS_EXTENTION}
