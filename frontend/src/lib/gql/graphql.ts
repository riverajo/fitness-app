/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: any; output: any; }
};

export type AuthPayload = {
  __typename?: 'AuthPayload';
  message: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  user?: Maybe<User>;
};

export type CreateUniqueExerciseInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
};

export type CreateWorkoutLogInput = {
  endTime: Scalars['Time']['input'];
  exerciseLogs: Array<ExerciseLogInput>;
  generalNotes?: InputMaybe<Scalars['String']['input']>;
  locationName?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
  startTime: Scalars['Time']['input'];
};

export type ExerciseLog = {
  __typename?: 'ExerciseLog';
  notes?: Maybe<Scalars['String']['output']>;
  sets: Array<Set>;
  uniqueExerciseId: Scalars['ID']['output'];
};

export type ExerciseLogInput = {
  notes?: InputMaybe<Scalars['String']['input']>;
  sets: Array<SetInput>;
  uniqueExerciseId: Scalars['ID']['input'];
};

export type LoginInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createUniqueExercise: UniqueExercise;
  createWorkoutLog: WorkoutLog;
  login: AuthPayload;
  logout: AuthPayload;
  register: AuthPayload;
  updateUser: AuthPayload;
};


export type MutationCreateUniqueExerciseArgs = {
  input: CreateUniqueExerciseInput;
};


export type MutationCreateWorkoutLogArgs = {
  input: CreateWorkoutLogInput;
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationUpdateUserArgs = {
  input: UpdateUserInput;
};

export type Query = {
  __typename?: 'Query';
  getWorkoutLog?: Maybe<WorkoutLog>;
  listWorkoutLogs: Array<WorkoutLog>;
  me?: Maybe<User>;
  searchExercises: Array<UniqueExercise>;
};


export type QueryGetWorkoutLogArgs = {
  id: Scalars['ID']['input'];
};


export type QuerySearchExercisesArgs = {
  query: Scalars['String']['input'];
};

export type RegisterInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Set = {
  __typename?: 'Set';
  order: Scalars['Int']['output'];
  reps: Scalars['Int']['output'];
  rpe?: Maybe<Scalars['Int']['output']>;
  toFailure?: Maybe<Scalars['Boolean']['output']>;
  weight: Scalars['Float']['output'];
};

export type SetInput = {
  order: Scalars['Int']['input'];
  reps: Scalars['Int']['input'];
  rpe?: InputMaybe<Scalars['Int']['input']>;
  toFailure?: InputMaybe<Scalars['Boolean']['input']>;
  unit: WeightUnit;
  weight: Scalars['Float']['input'];
};

export type UniqueExercise = {
  __typename?: 'UniqueExercise';
  description?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  isCustom: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
};

export type UpdateUserInput = {
  currentPassword: Scalars['String']['input'];
  email?: InputMaybe<Scalars['String']['input']>;
  newPassword?: InputMaybe<Scalars['String']['input']>;
  preferredUnit?: InputMaybe<Scalars['String']['input']>;
};

export type User = {
  __typename?: 'User';
  email: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  preferredUnit: Scalars['String']['output'];
};

export enum WeightUnit {
  Kilograms = 'KILOGRAMS',
  Pounds = 'POUNDS'
}

export type WorkoutLog = {
  __typename?: 'WorkoutLog';
  endTime: Scalars['Time']['output'];
  exerciseLogs: Array<ExerciseLog>;
  generalNotes?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  locationName?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  startTime: Scalars['Time']['output'];
};

export type LoginMutationVariables = Exact<{
  input: LoginInput;
}>;


export type LoginMutation = { __typename?: 'Mutation', login: { __typename?: 'AuthPayload', success: boolean, message: string, user?: { __typename?: 'User', id: string, email: string } | null } };

export type ListWorkoutLogsQueryVariables = Exact<{ [key: string]: never; }>;


export type ListWorkoutLogsQuery = { __typename?: 'Query', listWorkoutLogs: Array<{ __typename?: 'WorkoutLog', id: string, name: string, startTime: any, endTime: any, locationName?: string | null, exerciseLogs: Array<{ __typename?: 'ExerciseLog', uniqueExerciseId: string }> }> };

export type CreateUniqueExerciseMutationVariables = Exact<{
  input: CreateUniqueExerciseInput;
}>;


export type CreateUniqueExerciseMutation = { __typename?: 'Mutation', createUniqueExercise: { __typename?: 'UniqueExercise', id: string, name: string, isCustom: boolean } };

export type RegisterMutationVariables = Exact<{
  input: RegisterInput;
}>;


export type RegisterMutation = { __typename?: 'Mutation', register: { __typename?: 'AuthPayload', success: boolean, message: string, user?: { __typename?: 'User', id: string, email: string } | null } };

export type GetWorkoutLogQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetWorkoutLogQuery = { __typename?: 'Query', getWorkoutLog?: { __typename?: 'WorkoutLog', id: string, name: string, startTime: any, endTime: any, locationName?: string | null, generalNotes?: string | null, exerciseLogs: Array<{ __typename?: 'ExerciseLog', uniqueExerciseId: string, notes?: string | null, sets: Array<{ __typename?: 'Set', reps: number, weight: number, rpe?: number | null, toFailure?: boolean | null }> }> } | null };

export type CreateWorkoutLogMutationVariables = Exact<{
  input: CreateWorkoutLogInput;
}>;


export type CreateWorkoutLogMutation = { __typename?: 'Mutation', createWorkoutLog: { __typename?: 'WorkoutLog', id: string } };

export type SearchExercisesQueryVariables = Exact<{
  query: Scalars['String']['input'];
}>;


export type SearchExercisesQuery = { __typename?: 'Query', searchExercises: Array<{ __typename?: 'UniqueExercise', id: string, name: string, description?: string | null, isCustom: boolean }> };


export const LoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Login"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"LoginInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}}]}}]}}]}}]} as unknown as DocumentNode<LoginMutation, LoginMutationVariables>;
export const ListWorkoutLogsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ListWorkoutLogs"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"listWorkoutLogs"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"startTime"}},{"kind":"Field","name":{"kind":"Name","value":"endTime"}},{"kind":"Field","name":{"kind":"Name","value":"locationName"}},{"kind":"Field","name":{"kind":"Name","value":"exerciseLogs"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"uniqueExerciseId"}}]}}]}}]}}]} as unknown as DocumentNode<ListWorkoutLogsQuery, ListWorkoutLogsQueryVariables>;
export const CreateUniqueExerciseDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateUniqueExercise"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateUniqueExerciseInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createUniqueExercise"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"isCustom"}}]}}]}}]} as unknown as DocumentNode<CreateUniqueExerciseMutation, CreateUniqueExerciseMutationVariables>;
export const RegisterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Register"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"RegisterInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"register"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}}]}}]}}]}}]} as unknown as DocumentNode<RegisterMutation, RegisterMutationVariables>;
export const GetWorkoutLogDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetWorkoutLog"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"getWorkoutLog"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"startTime"}},{"kind":"Field","name":{"kind":"Name","value":"endTime"}},{"kind":"Field","name":{"kind":"Name","value":"locationName"}},{"kind":"Field","name":{"kind":"Name","value":"generalNotes"}},{"kind":"Field","name":{"kind":"Name","value":"exerciseLogs"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"uniqueExerciseId"}},{"kind":"Field","name":{"kind":"Name","value":"notes"}},{"kind":"Field","name":{"kind":"Name","value":"sets"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"reps"}},{"kind":"Field","name":{"kind":"Name","value":"weight"}},{"kind":"Field","name":{"kind":"Name","value":"rpe"}},{"kind":"Field","name":{"kind":"Name","value":"toFailure"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetWorkoutLogQuery, GetWorkoutLogQueryVariables>;
export const CreateWorkoutLogDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateWorkoutLog"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateWorkoutLogInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createWorkoutLog"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<CreateWorkoutLogMutation, CreateWorkoutLogMutationVariables>;
export const SearchExercisesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"SearchExercises"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"query"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"searchExercises"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"query"},"value":{"kind":"Variable","name":{"kind":"Name","value":"query"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"isCustom"}}]}}]}}]} as unknown as DocumentNode<SearchExercisesQuery, SearchExercisesQueryVariables>;