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
  createWorkoutLog: WorkoutLog;
  login: AuthPayload;
  logout: AuthPayload;
  register: AuthPayload;
  updateUser: AuthPayload;
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
};


export type QueryGetWorkoutLogArgs = {
  id: Scalars['ID']['input'];
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

export type ListWorkoutsQueryVariables = Exact<{ [key: string]: never; }>;


export type ListWorkoutsQuery = { __typename?: 'Query', listWorkoutLogs: Array<{ __typename?: 'WorkoutLog', id: string, name: string }> };


export const ListWorkoutsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ListWorkouts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"listWorkoutLogs"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]} as unknown as DocumentNode<ListWorkoutsQuery, ListWorkoutsQueryVariables>;