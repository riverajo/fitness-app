import { gql, type Client } from '@urql/svelte';

export interface WorkoutState {
	id: string | null;
	name: string;
	locationName: string;
	generalNotes: string;
	startTime: string;
	endTime: string;
	exerciseLogs: {
		uniqueExerciseId: string;
		name: string;
		sets: {
			reps: number;
			weight: number;
			unit: string;
			rpe?: number | null;
			toFailure?: boolean | null;
			order: number;
		}[];
		notes: string;
	}[];
}

const initialState: WorkoutState = {
	id: null,
	name: '',
	locationName: '',
	generalNotes: '',
	startTime: getISOString(), // This will be overwritten on reset, but good for type safety
	endTime: getISOString(),
	exerciseLogs: []
};

function getISOString(date?: Date | string | number) {
	if (date) {
		return new Date(date).toISOString();
	}
	return new Date().toISOString();
}

function validateWorkoutAge(startTime: string) {
	const workoutTime = new Date(startTime).getTime();
	const now = Date.now();
	const msIn24Hours = 24 * 60 * 60 * 1000;

	if (now - workoutTime > msIn24Hours) {
		throw new Error('Workout is too old to edit (limit: 24 hours)');
	}
}

class Workout {
	#state = $state<WorkoutState>(JSON.parse(JSON.stringify(initialState)));

	// Public accessor for the state object
	// We expose the raw state object so deep binding works (e.g. bind:value={workoutStore.state.name})
	get state() {
		return this.#state;
	}

	set state(newState: WorkoutState) {
		this.#state = newState;
	}

	reset() {
		this.#state = {
			...JSON.parse(JSON.stringify(initialState)),
			startTime: getISOString(),
			endTime: getISOString()
		};
	}

	update(fn: (state: WorkoutState) => WorkoutState) {
		this.#state = fn(this.#state);
	}

	async loadWorkoutForEditing(client: Client, workoutId: string) {
		const GET_WORKOUT_FOR_EDIT = gql`
			query GetWorkoutForEdit($id: ID!) {
				getWorkoutLog(id: $id) {
					id
					name
					locationName
					generalNotes
					startTime
					endTime
					exerciseLogs {
						uniqueExercise {
							id
							name
						}
						sets {
							reps
							weight
							rpe
							toFailure
							order
						}
						notes
					}
				}
			}
		`;

		const result = await client.query(GET_WORKOUT_FOR_EDIT, { id: workoutId }).toPromise();

		if (result.error) {
			throw new Error(result.error.message);
		}

		if (!result.data?.getWorkoutLog) {
			throw new Error('Workout not found');
		}

		const workout = result.data.getWorkoutLog;

		// 24-hour validation

		// 24-hour validation
		validateWorkoutAge(workout.startTime);

		// Map response to store state
		this.#state = {
			id: workout.id,
			name: workout.name,
			locationName: workout.locationName || '',
			generalNotes: workout.generalNotes || '',
			startTime: workout.startTime,
			endTime: workout.endTime,
			exerciseLogs: workout.exerciseLogs.map(
				(log: {
					uniqueExercise: { id: string; name: string };
					sets: {
						reps: number;
						weight: number;
						rpe?: number | null;
						toFailure?: boolean | null;
						order: number;
					}[];
					notes?: string;
				}) => ({
					uniqueExerciseId: log.uniqueExercise.id,
					name: log.uniqueExercise.name,
					sets: log.sets.map(
						(s: {
							reps: number;
							weight: number;
							rpe?: number | null;
							toFailure?: boolean | null;
							order: number;
						}) => ({
							reps: s.reps,
							weight: s.weight,
							unit: 'KILOGRAMS', // Backend is always KGS
							rpe: s.rpe,
							toFailure: s.toFailure,
							order: s.order
						})
					),
					notes: log.notes || ''
				})
			)
		};
	}
}

export const workoutStore = new Workout();
export { Workout };
