<script lang="ts">
    import { gql, getContextClient } from '@urql/svelte';
    import { goto } from '$app/navigation';

    const client = getContextClient();

    let name = '';
    let locationName = '';
    let generalNotes = '';
    let error = '';
    let searchQuery = '';
    let searchResults: any[] = [];
    let selectedExercises: any[] = [];

    // Current exercise being added
    let currentExercise: any = null;
    let currentSets: any[] = [];
    
    // Set input fields
    let reps = 0;
    let weight = 0;
    let rpe = 0;
    let toFailure = false;

    const searchExercisesQuery = gql`
        query SearchExercises($query: String!) {
            searchExercises(query: $query) {
                id
                name
                description
                isCustom
            }
        }
    `;

    const createWorkoutLogMutation = gql`
        mutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {
            createWorkoutLog(input: $input) {
                id
                name
            }
        }
    `;

    async function handleSearch() {
        if (!searchQuery) return;
        const result = await client.query(searchExercisesQuery, { query: searchQuery }).toPromise();
        if (result.data) {
            searchResults = result.data.searchExercises;
        }
    }

    function selectExercise(exercise: any) {
        currentExercise = exercise;
        currentSets = [];
        searchResults = [];
        searchQuery = '';
    }

    function addSet() {
        currentSets = [...currentSets, {
            reps,
            weight,
            rpe,
            toFailure,
            order: currentSets.length + 1,
            unit: 'KILOGRAMS' // Defaulting to KGS for now as per schema requirement
        }];
        // Reset fields for next set
        reps = 0;
        weight = 0;
        rpe = 0;
        toFailure = false;
    }

    function confirmExercise() {
        if (currentExercise && currentSets.length > 0) {
            selectedExercises = [...selectedExercises, {
                uniqueExerciseId: currentExercise.id,
                name: currentExercise.name, // For display
                sets: currentSets,
                notes: ''
            }];
            currentExercise = null;
            currentSets = [];
        }
    }

    async function handleSubmit() {
        error = '';
        
        const exerciseLogs = selectedExercises.map(ex => ({
            uniqueExerciseId: ex.uniqueExerciseId,
            sets: ex.sets.map((s: any) => ({
                reps: s.reps,
                weight: s.weight,
                unit: s.unit,
                rpe: s.rpe,
                toFailure: s.toFailure,
                order: s.order
            })),
            notes: ex.notes
        }));

        const input = {
            name,
            startTime: new Date().toISOString(),
            endTime: new Date().toISOString(), // Just using same time for now, ideally would be actual end time
            locationName,
            generalNotes,
            exerciseLogs
        };

        const result = await client.mutation(createWorkoutLogMutation, { input }).toPromise();

        if (result.error) {
            error = result.error.message;
        } else {
            goto('/dashboard');
        }
    }
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
    <h1 class="text-3xl font-bold tracking-tight text-gray-900">Log Workout</h1>

    <div class="mt-6 space-y-6 max-w-2xl">
        {#if error}
            <div class="rounded-md bg-red-50 p-4">
                <p class="text-sm text-red-700">{error}</p>
            </div>
        {/if}

        <!-- Workout Details -->
        <div class="space-y-4 rounded-lg border border-gray-200 p-4">
            <h2 class="text-lg font-semibold">Details</h2>
            <div>
                <label for="name" class="block text-sm font-medium text-gray-700">Workout Name</label>
                <input type="text" id="name" bind:value={name} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" placeholder="e.g. Morning Lift">
            </div>
            <div>
                <label for="location" class="block text-sm font-medium text-gray-700">Location</label>
                <input type="text" id="location" bind:value={locationName} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
            </div>
            <div>
                <label for="notes" class="block text-sm font-medium text-gray-700">Notes</label>
                <textarea id="notes" bind:value={generalNotes} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"></textarea>
            </div>
        </div>

        <!-- Added Exercises List -->
        {#if selectedExercises.length > 0}
            <div class="space-y-4">
                <h2 class="text-lg font-semibold">Exercises</h2>
                {#each selectedExercises as exercise, i}
                    <div class="rounded-lg border border-gray-200 p-4 bg-gray-50">
                        <h3 class="font-medium">{exercise.name}</h3>
                        <div class="mt-2 text-sm text-gray-600">
                            {#each exercise.sets as set}
                                <div>Set {set.order}: {set.reps} reps @ {set.weight}kg</div>
                            {/each}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}

        <!-- Add Exercise Section -->
        <div class="rounded-lg border border-gray-200 p-4">
            <h2 class="text-lg font-semibold">Add Exercise</h2>
            
            {#if !currentExercise}
                <div class="mt-4 flex gap-2">
                    <input 
                        type="text" 
                        bind:value={searchQuery} 
                        placeholder="Search exercises..." 
                        class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        on:keydown={(e) => e.key === 'Enter' && handleSearch()}
                    >
                    <button on:click={handleSearch} class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500">Search</button>
                </div>

                {#if searchResults.length > 0}
                    <ul class="mt-4 divide-y divide-gray-200 rounded-md border border-gray-200">
                        {#each searchResults as exercise}
                            <li>
                                <button class="w-full px-4 py-2 text-left hover:bg-gray-50" on:click={() => selectExercise(exercise)}>
                                    <span class="font-medium">{exercise.name}</span>
                                    {#if exercise.description}
                                        <span class="ml-2 text-sm text-gray-500">- {exercise.description}</span>
                                    {/if}
                                </button>
                            </li>
                        {/each}
                    </ul>
                {/if}
            {:else}
                <div class="mt-4 space-y-4">
                    <div class="flex justify-between items-center">
                        <h3 class="font-medium text-indigo-600">{currentExercise.name}</h3>
                        <button on:click={() => currentExercise = null} class="text-sm text-gray-500 hover:text-gray-700">Change Exercise</button>
                    </div>

                    <!-- Sets List -->
                    {#if currentSets.length > 0}
                        <div class="space-y-2">
                            {#each currentSets as set}
                                <div class="text-sm">Set {set.order}: {set.reps} reps @ {set.weight}kg</div>
                            {/each}
                        </div>
                    {/if}

                    <!-- Add Set Form -->
                    <div class="grid grid-cols-2 gap-4 bg-gray-50 p-3 rounded-md">
                        <div>
                            <label for="reps" class="block text-xs font-medium text-gray-700">Reps</label>
                            <input type="number" id="reps" bind:value={reps} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-xs">
                        </div>
                        <div>
                            <label for="weight" class="block text-xs font-medium text-gray-700">Weight (kg)</label>
                            <input type="number" id="weight" bind:value={weight} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-xs">
                        </div>
                        <div>
                            <label for="rpe" class="block text-xs font-medium text-gray-700">RPE</label>
                            <input type="number" id="rpe" bind:value={rpe} class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-xs">
                        </div>
                        <div class="flex items-center mt-6">
                            <input type="checkbox" id="toFailure" bind:checked={toFailure} class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500">
                            <label for="toFailure" class="ml-2 block text-xs text-gray-900">To Failure</label>
                        </div>
                        <div class="col-span-2">
                            <button on:click={addSet} class="w-full rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">Add Set</button>
                        </div>
                    </div>

                    <button on:click={confirmExercise} class="w-full rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500">Done Adding Sets</button>
                </div>
            {/if}
        </div>

        <div class="flex justify-end gap-4 pt-4">
            <a href="/dashboard" class="rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">Cancel</a>
            <button on:click={handleSubmit} class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500">Save Workout</button>
        </div>
    </div>
</div>
