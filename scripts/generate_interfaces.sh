SCRIPT_DIRECTORY=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

ROOT_PROJECT_DIRECTORY="${SCRIPT_DIRECTORY}/.."
GENERATED_DECORATORS_EXTENTION=".gen.go"
GO_DECORATORS_TEMPLATE_DIRECTORY="${ROOT_PROJECT_DIRECTORY}/core/pkg/decorators/templates"

cd $ROOT_PROJECT_DIRECTORY

export GOPATH="/home/am/go"
export PATH="$PATH:$(go env GOPATH)/bin"

GO_WRAP_BINARY="${GOPATH}/bin/gowrap"

# user service
# user service: repositories
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories \
    -i UserRepository \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/ulo_logger.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/decorators/ulo_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories \
    -i UserRepository \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/opentelemetry_tracing.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/decorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -v metricPrefix=database_query \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories \
    -i UserRepository \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/prometheus_metrics.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/decorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# # user service: services
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services \
    -i UserService \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/ulo_logger.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services/decorators/ulo_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services \
    -i UserService \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/opentelemetry_tracing.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services/decorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -v metricPrefix=service_request \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services \
    -i UserService \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/prometheus_metrics.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services/decorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# # user service: controllers
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/user_service/internal/controllers \
    -i UserController \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/ulo_logger.go \
    -o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/controllers/decorators/ulo_logger${GENERATED_DECORATORS_EXTENTION}

# auth service
# auth service: repositories
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/repositories \
    -i AccountRepository \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/ulo_logger.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/repositories/decorators/ulo_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/repositories \
    -i AccountRepository \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/opentelemetry_tracing.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/repositories/decorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -v metricPrefix=database_query \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/repositories \
    -i AccountRepository \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/prometheus_metrics.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/repositories/decorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# auth service: services
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/services \
    -i AccountService \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/ulo_logger.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/services/decorators/ulo_logger${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/services \
    -i AccountService \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/opentelemetry_tracing.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/services/decorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

${GO_WRAP_BINARY} gen \
    -v metricPrefix=service_request \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/services \
    -i AccountService \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/prometheus_metrics.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/services/decorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

# auth service: controllers
${GO_WRAP_BINARY} gen \
    -p ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/controllers \
    -i AccountController \
    -t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/ulo_logger.go \
    -o ${ROOT_PROJECT_DIRECTORY}/auth_service/internal/controllers/decorators/ulo_logger${GENERATED_DECORATORS_EXTENTION}
