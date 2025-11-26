<script lang="ts">
    import { gql, getContextClient } from '@urql/svelte';
    import { Heading, Button, Card, Label, Input, Textarea, Checkbox, Alert, Listgroup, ListgroupItem } from 'flowbite-svelte';
    import { createEventDispatcher } from 'svelte';

    export let initialData: any = null;
    export let submitLabel: string = 'Save Workout';
    export let error: string = '';

    const dispatch = createEventDispatcher();
    const client = getContextClient();

    let name = initialData?.name || '';
    let locationName = initialData?.locationName || '';
    let generalNotes = initialData?.generalNotes || '';
    
    let searchQuery = '';
    let searchResults: any[] = [];
    
    // Initialize selected exercises from initialData if available
    let selectedExercises: any[] = initialData?.exerciseLogs?.map((log: any) => ({
        uniqueExerciseId: log.uniqueExercise.id,
        name: log.uniqueExercise.name,
        sets: log.sets.map((s: any) => ({ ...s })),
        notes: log.notes || ''
    })) || [];

    // Current exercise being added
    let currentExercise: any = null;
    let currentSets: any[] = [];
    
    // Set input fields
    let reps = 0;
    let weight = 0;
    let rpe = 0;
    let toFailure = false;

    const searchExercisesQuery = gql`
        query SearchUniqueExercises($query: String) {
            uniqueExercises(query: $query, limit: 10) {
                id
                name
                description
                isCustom
            }
        }
    `;

    async function handleSearch() {
        if (!searchQuery) return;
        const result = await client.query(searchExercisesQuery, { query: searchQuery }).toPromise();
        if (result.data) {
            searchResults = result.data.uniqueExercises;
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
            unit: 'KILOGRAMS'
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
                name: currentExercise.name,
                sets: currentSets,
                notes: ''
            }];
            currentExercise = null;
            currentSets = [];
        }
    }

    function removeExercise(index: number) {
        selectedExercises = selectedExercises.filter((_, i) => i !== index);
    }

    function handleSubmit() {
        const exerciseLogs = selectedExercises.map(ex => ({
            uniqueExerciseId: ex.uniqueExerciseId,
            sets: ex.sets.map((s: any, i: number) => ({
                reps: s.reps,
                weight: s.weight,
                unit: s.unit || 'KILOGRAMS',
                rpe: s.rpe,
                toFailure: s.toFailure,
                order: i + 1 // Ensure order is correct
            })),
            notes: ex.notes
        }));

        const formData = {
            name,
            locationName,
            generalNotes,
            exerciseLogs
        };

        dispatch('submit', formData);
    }
</script>

<div class="space-y-6 max-w-2xl">
    {#if error}
        <Alert color="red">{error}</Alert>
    {/if}

    <!-- Workout Details -->
    <Card class="w-full max-w-none">
        <Heading tag="h2" class="mb-4 text-lg">Details</Heading>
        <div class="space-y-4">
            <div>
                <Label for="name" class="mb-2">Workout Name</Label>
                <Input type="text" id="name" bind:value={name} placeholder="e.g. Morning Lift" />
            </div>
            <div>
                <Label for="location" class="mb-2">Location</Label>
                <Input type="text" id="location" bind:value={locationName} />
            </div>
            <div>
                <Label for="notes" class="mb-2">Notes</Label>
                <Textarea id="notes" bind:value={generalNotes} />
            </div>
        </div>
    </Card>

    <!-- Added Exercises List -->
    {#if selectedExercises.length > 0}
        <div class="space-y-4">
            <Heading tag="h2" class="text-lg">Exercises</Heading>
            {#each selectedExercises as exercise, i}
                <Card class="w-full max-w-none bg-gray-50 dark:bg-gray-700">
                    <div class="flex justify-between items-start">
                        <h3 class="font-medium text-gray-900 dark:text-white">{exercise.name}</h3>
                        <Button color="red" size="xs" onclick={() => removeExercise(i)}>Remove</Button>
                    </div>
                    <div class="mt-2 text-sm text-gray-600 dark:text-gray-300">
                        {#each exercise.sets as set}
                            <div>Set {set.order}: {set.reps} reps @ {set.weight}kg</div>
                        {/each}
                    </div>
                </Card>
            {/each}
        </div>
    {/if}

    <!-- Add Exercise Section -->
    <Card class="w-full max-w-none">
        <Heading tag="h2" class="mb-4 text-lg">Add Exercise</Heading>
        
        {#if !currentExercise}
            <div class="mt-4 flex gap-2">
                <div class="flex-1">
                    <Input 
                        type="text" 
                        bind:value={searchQuery} 
                        placeholder="Search exercises..." 
                        onkeydown={(e) => e.key === 'Enter' && handleSearch()}
                    />
                </div>
                <Button onclick={handleSearch} color="blue">Search</Button>
            </div>

            {#if searchResults.length > 0}
                <Listgroup class="mt-4">
                    {#each searchResults as exercise}
                        <ListgroupItem class="p-0">
                            <button 
                                type="button" 
                                class="w-full text-left px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors duration-200 flex items-center" 
                                onclick={() => selectExercise(exercise)}
                            >
                                <span class="font-medium">{exercise.name}</span>
                                {#if exercise.description}
                                    <span class="ml-2 text-sm text-gray-500">- {exercise.description}</span>
                                {/if}
                            </button>
                        </ListgroupItem>
                    {/each}
                </Listgroup>
            {/if}
        {:else}
            <div class="mt-4 space-y-4">
                <div class="flex justify-between items-center">
                    <h3 class="font-medium text-indigo-600 dark:text-indigo-400">{currentExercise.name}</h3>
                    <Button color="light" size="xs" onclick={() => currentExercise = null}>Change Exercise</Button>
                </div>

                <!-- Sets List -->
                {#if currentSets.length > 0}
                    <div class="space-y-2">
                        {#each currentSets as set}
                            <div class="text-sm text-gray-700 dark:text-gray-300">Set {set.order}: {set.reps} reps @ {set.weight}kg</div>
                        {/each}
                    </div>
                {/if}

                <!-- Add Set Form -->
                <div class="grid grid-cols-2 gap-4 bg-gray-50 p-3 rounded-md dark:bg-gray-700">
                    <div>
                        <Label for="reps" class="mb-1 text-xs">Reps</Label>
                        <Input type="number" id="reps" bind:value={reps} size="sm" />
                    </div>
                    <div>
                        <Label for="weight" class="mb-1 text-xs">Weight (kg)</Label>
                        <Input type="number" id="weight" bind:value={weight} size="sm" />
                    </div>
                    <div>
                        <Label for="rpe" class="mb-1 text-xs">RPE</Label>
                        <Input type="number" id="rpe" bind:value={rpe} size="sm" />
                    </div>
                    <div class="flex items-center mt-6">
                        <Checkbox id="toFailure" bind:checked={toFailure}>To Failure</Checkbox>
                    </div>
                    <div class="col-span-2">
                        <Button onclick={addSet} color="light" class="w-full">Add Set</Button>
                    </div>
                </div>

                <Button onclick={confirmExercise} class="w-full">Done Adding Sets</Button>
            </div>
        {/if}
    </Card>

    <div class="flex justify-end gap-4 pt-4">
        <Button color="light" href="/dashboard">Cancel</Button>
        <Button onclick={handleSubmit}>{submitLabel}</Button>
    </div>
</div>
