<script lang="ts">
	import { Input } from 'flowbite-svelte';

	interface Props {
		value?: number;
		size?: 'sm' | 'md' | 'lg';
	}

	let { value = $bindable(0), size = 'sm' }: Props = $props();

	let lbs = $state<number | undefined>();
	let kgs = $state<number | undefined>(value === 0 ? undefined : value);

	// Constant for conversion
	const LBS_TO_KG = 0.45359237;

	// Track the last value we emitted to detecting external changes
	let lastEmittedValue = $state(value);

	// Watch for external changes to value (e.g. loading from DB or reset form)
	$effect(() => {
		if (value !== lastEmittedValue) {
			// If the value changes from outside, we reset our local inputs
			// Default to all KGs since we don't store the split
			kgs = value === 0 ? undefined : value;
			lbs = undefined;
			lastEmittedValue = value;
		}
	});

	function updateValue() {
		const l = lbs ?? 0;
		const k = kgs ?? 0;
		let total = l * LBS_TO_KG + k;
		// Round to 2 decimal places
		total = Math.round(total * 100) / 100;

		lastEmittedValue = total;
		value = total;
	}
</script>

<div class="flex items-center gap-2">
	<div class="flex flex-col">
		<div class="flex items-center gap-1">
			<Input
				type="number"
				{size}
				class="w-20"
				bind:value={lbs}
				oninput={updateValue}
				placeholder="0"
				data-testid="weight-input-lbs"
			/>
			<span class="text-xs text-gray-500">lbs</span>
		</div>
	</div>
	<span class="text-gray-400">+</span>
	<div class="flex flex-col">
		<div class="flex items-center gap-1">
			<Input
				type="number"
				{size}
				class="w-20"
				bind:value={kgs}
				oninput={updateValue}
				placeholder="0"
				data-testid="weight-input-kgs"
			/>
			<span class="text-xs text-gray-500">kgs</span>
		</div>
	</div>
</div>
