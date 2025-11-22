/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n        mutation Login($input: LoginInput!) {\n            login(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    ": typeof types.LoginDocument,
    "\n        query ListWorkoutLogs {\n            listWorkoutLogs {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                exerciseLogs {\n                    uniqueExerciseId\n                }\n            }\n        }\n    ": typeof types.ListWorkoutLogsDocument,
    "\n        mutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n            createUniqueExercise(input: $input) {\n                id\n                name\n                isCustom\n            }\n        }\n    ": typeof types.CreateUniqueExerciseDocument,
    "\n        mutation Register($input: RegisterInput!) {\n            register(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    ": typeof types.RegisterDocument,
    "\n        query GetWorkoutLog($id: ID!) {\n            getWorkoutLog(id: $id) {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                generalNotes\n                exerciseLogs {\n                    uniqueExerciseId\n                    notes\n                    sets {\n                        reps\n                        weight\n                        rpe\n                        toFailure\n                    }\n                }\n            }\n        }\n    ": typeof types.GetWorkoutLogDocument,
    "\n        mutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n            createWorkoutLog(input: $input) {\n                id\n            }\n        }\n    ": typeof types.CreateWorkoutLogDocument,
    "\n        query SearchExercises($query: String!) {\n            searchExercises(query: $query) {\n                id\n                name\n                description\n                isCustom\n            }\n        }\n    ": typeof types.SearchExercisesDocument,
};
const documents: Documents = {
    "\n        mutation Login($input: LoginInput!) {\n            login(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    ": types.LoginDocument,
    "\n        query ListWorkoutLogs {\n            listWorkoutLogs {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                exerciseLogs {\n                    uniqueExerciseId\n                }\n            }\n        }\n    ": types.ListWorkoutLogsDocument,
    "\n        mutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n            createUniqueExercise(input: $input) {\n                id\n                name\n                isCustom\n            }\n        }\n    ": types.CreateUniqueExerciseDocument,
    "\n        mutation Register($input: RegisterInput!) {\n            register(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    ": types.RegisterDocument,
    "\n        query GetWorkoutLog($id: ID!) {\n            getWorkoutLog(id: $id) {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                generalNotes\n                exerciseLogs {\n                    uniqueExerciseId\n                    notes\n                    sets {\n                        reps\n                        weight\n                        rpe\n                        toFailure\n                    }\n                }\n            }\n        }\n    ": types.GetWorkoutLogDocument,
    "\n        mutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n            createWorkoutLog(input: $input) {\n                id\n            }\n        }\n    ": types.CreateWorkoutLogDocument,
    "\n        query SearchExercises($query: String!) {\n            searchExercises(query: $query) {\n                id\n                name\n                description\n                isCustom\n            }\n        }\n    ": types.SearchExercisesDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        mutation Login($input: LoginInput!) {\n            login(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    "): (typeof documents)["\n        mutation Login($input: LoginInput!) {\n            login(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        query ListWorkoutLogs {\n            listWorkoutLogs {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                exerciseLogs {\n                    uniqueExerciseId\n                }\n            }\n        }\n    "): (typeof documents)["\n        query ListWorkoutLogs {\n            listWorkoutLogs {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                exerciseLogs {\n                    uniqueExerciseId\n                }\n            }\n        }\n    "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        mutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n            createUniqueExercise(input: $input) {\n                id\n                name\n                isCustom\n            }\n        }\n    "): (typeof documents)["\n        mutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n            createUniqueExercise(input: $input) {\n                id\n                name\n                isCustom\n            }\n        }\n    "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        mutation Register($input: RegisterInput!) {\n            register(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    "): (typeof documents)["\n        mutation Register($input: RegisterInput!) {\n            register(input: $input) {\n                success\n                message\n                user {\n                    id\n                    email\n                }\n            }\n        }\n    "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        query GetWorkoutLog($id: ID!) {\n            getWorkoutLog(id: $id) {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                generalNotes\n                exerciseLogs {\n                    uniqueExerciseId\n                    notes\n                    sets {\n                        reps\n                        weight\n                        rpe\n                        toFailure\n                    }\n                }\n            }\n        }\n    "): (typeof documents)["\n        query GetWorkoutLog($id: ID!) {\n            getWorkoutLog(id: $id) {\n                id\n                name\n                startTime\n                endTime\n                locationName\n                generalNotes\n                exerciseLogs {\n                    uniqueExerciseId\n                    notes\n                    sets {\n                        reps\n                        weight\n                        rpe\n                        toFailure\n                    }\n                }\n            }\n        }\n    "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        mutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n            createWorkoutLog(input: $input) {\n                id\n            }\n        }\n    "): (typeof documents)["\n        mutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n            createWorkoutLog(input: $input) {\n                id\n            }\n        }\n    "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n        query SearchExercises($query: String!) {\n            searchExercises(query: $query) {\n                id\n                name\n                description\n                isCustom\n            }\n        }\n    "): (typeof documents)["\n        query SearchExercises($query: String!) {\n            searchExercises(query: $query) {\n                id\n                name\n                description\n                isCustom\n            }\n        }\n    "];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;