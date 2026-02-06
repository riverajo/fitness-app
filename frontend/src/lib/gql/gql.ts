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
    "\n\t\tquery SearchUniqueExercises($query: String) {\n\t\t\tuniqueExercises(query: $query, limit: 10) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t": typeof types.SearchUniqueExercisesDocument,
    "\n\t\tquery Me {\n\t\t\tme {\n\t\t\t\tid\n\t\t\t\temail\n\t\t\t}\n\t\t}\n\t": typeof types.MeDocument,
    "\n\t\tmutation Logout {\n\t\t\tlogout {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t}\n\t\t}\n\t": typeof types.LogoutDocument,
    "\n\t\tmutation Login($input: LoginInput!) {\n\t\t\tlogin(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t": typeof types.LoginDocument,
    "\n\t\t\tquery ListWorkoutLogs($limit: Int, $offset: Int) {\n\t\t\t\tlistWorkoutLogs(limit: $limit, offset: $offset) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": typeof types.ListWorkoutLogsDocument,
    "\n\t\t\tquery Me {\n\t\t\t\tme {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t}\n\t\t": typeof types.MeDocument,
    "\n\t\t\t\tquery UniqueExercises($query: String, $limit: Int, $offset: Int) {\n\t\t\t\t\tuniqueExercises(query: $query, limit: $limit, offset: $offset) {\n\t\t\t\t\t\tid\n\t\t\t\t\t\tname\n\t\t\t\t\t\tdescription\n\t\t\t\t\t\tisCustom\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t": typeof types.UniqueExercisesDocument,
    "\n\t\tmutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n\t\t\tcreateUniqueExercise(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t": typeof types.CreateUniqueExerciseDocument,
    "\n\t\tmutation Register($input: RegisterInput!) {\n\t\t\tregister(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t": typeof types.RegisterDocument,
    "\n\t\t\tquery GetWorkoutLog($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": typeof types.GetWorkoutLogDocument,
    "\n\t\tmutation UpdateWorkoutLog($input: UpdateWorkoutLogInput!) {\n\t\t\tupdateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t": typeof types.UpdateWorkoutLogDocument,
    "\n\t\tmutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n\t\t\tcreateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t": typeof types.CreateWorkoutLogDocument,
    "\n\t\t\tquery GetWorkoutForEdit($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t\torder\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": typeof types.GetWorkoutForEditDocument,
};
const documents: Documents = {
    "\n\t\tquery SearchUniqueExercises($query: String) {\n\t\t\tuniqueExercises(query: $query, limit: 10) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t": types.SearchUniqueExercisesDocument,
    "\n\t\tquery Me {\n\t\t\tme {\n\t\t\t\tid\n\t\t\t\temail\n\t\t\t}\n\t\t}\n\t": types.MeDocument,
    "\n\t\tmutation Logout {\n\t\t\tlogout {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t}\n\t\t}\n\t": types.LogoutDocument,
    "\n\t\tmutation Login($input: LoginInput!) {\n\t\t\tlogin(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t": types.LoginDocument,
    "\n\t\t\tquery ListWorkoutLogs($limit: Int, $offset: Int) {\n\t\t\t\tlistWorkoutLogs(limit: $limit, offset: $offset) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": types.ListWorkoutLogsDocument,
    "\n\t\t\tquery Me {\n\t\t\t\tme {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t}\n\t\t": types.MeDocument,
    "\n\t\t\t\tquery UniqueExercises($query: String, $limit: Int, $offset: Int) {\n\t\t\t\t\tuniqueExercises(query: $query, limit: $limit, offset: $offset) {\n\t\t\t\t\t\tid\n\t\t\t\t\t\tname\n\t\t\t\t\t\tdescription\n\t\t\t\t\t\tisCustom\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t": types.UniqueExercisesDocument,
    "\n\t\tmutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n\t\t\tcreateUniqueExercise(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t": types.CreateUniqueExerciseDocument,
    "\n\t\tmutation Register($input: RegisterInput!) {\n\t\t\tregister(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t": types.RegisterDocument,
    "\n\t\t\tquery GetWorkoutLog($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": types.GetWorkoutLogDocument,
    "\n\t\tmutation UpdateWorkoutLog($input: UpdateWorkoutLogInput!) {\n\t\t\tupdateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t": types.UpdateWorkoutLogDocument,
    "\n\t\tmutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n\t\t\tcreateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t": types.CreateWorkoutLogDocument,
    "\n\t\t\tquery GetWorkoutForEdit($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t\torder\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": types.GetWorkoutForEditDocument,
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
export function graphql(source: "\n\t\tquery SearchUniqueExercises($query: String) {\n\t\t\tuniqueExercises(query: $query, limit: 10) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tquery SearchUniqueExercises($query: String) {\n\t\t\tuniqueExercises(query: $query, limit: 10) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tquery Me {\n\t\t\tme {\n\t\t\t\tid\n\t\t\t\temail\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tquery Me {\n\t\t\tme {\n\t\t\t\tid\n\t\t\t\temail\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation Logout {\n\t\t\tlogout {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation Logout {\n\t\t\tlogout {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation Login($input: LoginInput!) {\n\t\t\tlogin(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation Login($input: LoginInput!) {\n\t\t\tlogin(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery ListWorkoutLogs($limit: Int, $offset: Int) {\n\t\t\t\tlistWorkoutLogs(limit: $limit, offset: $offset) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery ListWorkoutLogs($limit: Int, $offset: Int) {\n\t\t\t\tlistWorkoutLogs(limit: $limit, offset: $offset) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery Me {\n\t\t\t\tme {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery Me {\n\t\t\t\tme {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\t\tquery UniqueExercises($query: String, $limit: Int, $offset: Int) {\n\t\t\t\t\tuniqueExercises(query: $query, limit: $limit, offset: $offset) {\n\t\t\t\t\t\tid\n\t\t\t\t\t\tname\n\t\t\t\t\t\tdescription\n\t\t\t\t\t\tisCustom\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t"): (typeof documents)["\n\t\t\t\tquery UniqueExercises($query: String, $limit: Int, $offset: Int) {\n\t\t\t\t\tuniqueExercises(query: $query, limit: $limit, offset: $offset) {\n\t\t\t\t\t\tid\n\t\t\t\t\t\tname\n\t\t\t\t\t\tdescription\n\t\t\t\t\t\tisCustom\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n\t\t\tcreateUniqueExercise(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {\n\t\t\tcreateUniqueExercise(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tdescription\n\t\t\t\tisCustom\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation Register($input: RegisterInput!) {\n\t\t\tregister(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation Register($input: RegisterInput!) {\n\t\t\tregister(input: $input) {\n\t\t\t\tsuccess\n\t\t\t\tmessage\n\t\t\t\tuser {\n\t\t\t\t\tid\n\t\t\t\t\temail\n\t\t\t\t}\n\t\t\t\ttoken\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery GetWorkoutLog($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery GetWorkoutLog($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation UpdateWorkoutLog($input: UpdateWorkoutLogInput!) {\n\t\t\tupdateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation UpdateWorkoutLog($input: UpdateWorkoutLogInput!) {\n\t\t\tupdateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n\t\t\tcreateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {\n\t\t\tcreateWorkoutLog(input: $input) {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery GetWorkoutForEdit($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t\torder\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery GetWorkoutForEdit($id: ID!) {\n\t\t\t\tgetWorkoutLog(id: $id) {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tlocationName\n\t\t\t\t\tgeneralNotes\n\t\t\t\t\tstartTime\n\t\t\t\t\tendTime\n\t\t\t\t\texerciseLogs {\n\t\t\t\t\t\tuniqueExercise {\n\t\t\t\t\t\t\tid\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t}\n\t\t\t\t\t\tsets {\n\t\t\t\t\t\t\treps\n\t\t\t\t\t\t\tweight\n\t\t\t\t\t\t\trpe\n\t\t\t\t\t\t\ttoFailure\n\t\t\t\t\t\t\torder\n\t\t\t\t\t\t}\n\t\t\t\t\t\tnotes\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;