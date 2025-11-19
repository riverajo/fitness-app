import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
    overwrite: true,
    schema: '../backend/graph/schema.graphqls',
    documents: ['src/**/*.svelte', 'src/**/*.ts'],
    ignoreNoDocuments: true, // for better experience with the watcher
    generates: {
        './src/lib/gql/': {
            preset: 'client',
            config: {
                useTypeImports: true
            },
            plugins: []
        }
    }
};

export default config;
