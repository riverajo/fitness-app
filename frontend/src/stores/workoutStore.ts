import { writable } from 'svelte/store';
import { gql, type Client } from '@urql/svelte'; // Using type import for Client to avoid runtime issues if strict

// Define the store state interface
// We use a shape similar to CreateWorkoutLogInput but with more UI-friendly types if needed
export interface WorkoutState {
	id: string | null; // null for new workouts
	name: string;
	locationName: string;
	generalNotes: string;
	startTime: string; // ISO string
	endTime: string; // ISO string - Added for full document replacement
	exerciseLogs: {
		uniqueExerciseId: string;
		name: string; // for display
		sets: {
			reps: number;
			weight: number;
			unit: string; // 'KILOGRAMS' | 'POUNDS'
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
	startTime: new Date().toISOString(),
	endTime: new Date().toISOString(),
	exerciseLogs: []
};

function createWorkoutStore() {
	const { subscribe, set, update } = writable<WorkoutState>(
		JSON.parse(JSON.stringify(initialState))
	);

	return {
		subscribe,
		set,
		update,
		reset: () => set(JSON.parse(JSON.stringify(initialState))),

		// Hydrate the store from a specific workout ID
		loadWorkoutForEditing: async (client: Client, workoutId: string) => {
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
			const workoutTime = new Date(workout.startTime).getTime();
			const now = Date.now();
			const msIn24Hours = 24 * 60 * 60 * 1000;

			if (now - workoutTime > msIn24Hours) {
				throw new Error('Workout is too old to edit (limit: 24 hours)');
			}

			// Map response to store state
			const mappedState: WorkoutState = {
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
								unit: 'KILOGRAMS', // Backend is always KGS, we default to that.
								// Future: could infer preference from user settings if available
								rpe: s.rpe,
								toFailure: s.toFailure,
								order: s.order
							})
						),
						notes: log.notes || ''
					})
				)
			};

			set(mappedState);
		}
	};
}

export const workoutStore = createWorkoutStore();
