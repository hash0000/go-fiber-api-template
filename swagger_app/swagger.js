const swaggerJSDoc = require("swagger-jsdoc");
const fs = require("node:fs");
const join = require("node:path").join;

const swaggerDefinition = {
    openapi: "3.0.0",
    info: {
        title: "template",
        version: "0.0.1",
    },
    components: {
        securitySchemes: {
            ApiTokenAuth: {
                type: "apiKey",
                in: "header",
                name: "api-token",
            },
        },
    },
};

const options = {
    swaggerDefinition,
    apis: [join(process.cwd() + `/external/documentation/raw/**/*.yaml`)],
};

const swaggerSpec = swaggerJSDoc(options);

fs.writeFileSync(join(process.cwd() + `/external/documentation/swagger.yaml`), JSON.stringify(swaggerSpec, null, 2));
