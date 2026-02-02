import { describe, it, expect, beforeEach, vi } from 'vitest';
import { Workout } from './workout.svelte';
import type { Client } from '@urql/svelte';

describe('Workout Store', () => {
	let workout: Workout;

	beforeEach(() => {
		workout = new Workout();
	});

	it('initializes with default state', () => {
		expect(workout.state.id).toBe(null);
		expect(workout.state.exerciseLogs).toEqual([]);
		expect(workout.state.name).toBe('');
	});

	it('resets state correctly', () => {
		workout.state.name = 'Test Workout';
		workout.reset();
		expect(workout.state.name).toBe('');
		// Start/end times should be recent
		const now = new Date().getTime();
		const startTime = new Date(workout.state.startTime).getTime();
		expect(Math.abs(now - startTime)).toBeLessThan(1000); // within 1 second
	});

	describe('loadWorkoutForEditing', () => {
		it('populates state from query result', async () => {
			const mockClient = {
				query: vi.fn().mockReturnValue({
					toPromise: vi.fn().mockResolvedValue({
						data: {
							getWorkoutLog: {
								id: '123',
								name: 'Old Workout',
								startTime: new Date().toISOString(), // Recent enough
								endTime: new Date().toISOString(),
								locationName: 'Gym',
								generalNotes: 'Notes',
								exerciseLogs: [
									{
										uniqueExercise: { id: 'ex1', name: 'Bench' },
										sets: [{ reps: 10, weight: 100, order: 1 }],
										notes: 'Easy'
									}
								]
							}
						}
					})
				})
			} as unknown as Client;

			await workout.loadWorkoutForEditing(mockClient, '123');

			expect(workout.state.id).toBe('123');
			expect(workout.state.name).toBe('Old Workout');
			expect(workout.state.exerciseLogs[0].name).toBe('Bench');
			expect(workout.state.exerciseLogs[0].sets[0].unit).toBe('KILOGRAMS');
		});

		it('throws if workout is too old', async () => {
			const oldDate = new Date();
			oldDate.setDate(oldDate.getDate() - 2); // 2 days old

			const mockClient = {
				query: vi.fn().mockReturnValue({
					toPromise: vi.fn().mockResolvedValue({
						data: {
							getWorkoutLog: {
								id: '123',
								name: 'Old Workout',
								startTime: oldDate.toISOString(),
								endTime: oldDate.toISOString(),
								exerciseLogs: []
							}
						}
					})
				})
			} as unknown as Client;

			await expect(workout.loadWorkoutForEditing(mockClient, '123')).rejects.toThrow(
				'Workout is too old'
			);
		});
	});
});
