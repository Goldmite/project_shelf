<script lang="ts">
	import type { Book, BookProgress } from '$lib/types/responses/book';
	import { toReadableLanguage } from '$lib/types/enums/language';
	import type { PageProps } from './$types';
	import { bookOwners, groupMembers } from '$lib/stores/localMembers';
	import { goto } from '$app/navigation';
	import BackButton from '../../../../../../components/navigation/BackButton.svelte';
	import CoverWrapper from '../../../../../../components/shelf/CoverWrapper.svelte';
	import MembersList from '../../../../../../components/shelf/shared/MembersList.svelte';
	import PageSubheader from '../../../../../../components/PageSubheader.svelte';
	import { page } from '$app/state';
	import { user } from '$lib/stores/user';
	import { formatMinToHoursOrMin } from '$lib/helpers/formatMinTime';
	import { MINUTES_TO_SEE_STATS } from '$lib/constants';

	let { data }: PageProps = $props();
	const book: Book = data.book;
	const progress: BookProgress | undefined = data.progress;
	book.owned_by = $bookOwners;

	let readMore = $state(false);
	const descriptionPreview = 300;

	function calculateRemainingTime() {
		if (progress) {
			const pagesLeft = book.pages - progress.current_page;
			const ppm = progress.pages_read / (progress.time_read / 60);
			return formatMinToHoursOrMin(Math.round(pagesLeft / ppm));
		} else {
			return '_';
		}
	}
</script>

<div class="flex min-h-[80vh] flex-col gap-2 sm:gap-4 md:gap-6">
	<span class="my-3 flex flex-col items-baseline md:flex-row md:gap-4">
		<PageSubheader>
			<BackButton backUrl="/{page.params.from}/{page.params.sharedId}" />{book.title}
		</PageSubheader>
		<p class="inline font-light italic">
			by
			{#each book.authors as author, i}
				{author}{i + 1 < book.authors.length ? ', ' : ''}
			{/each}
		</p>
		{#if book.owned_by.includes($user?.id ?? 'NO_USER')}
			<div class="ml-auto flex flex-row gap-2 sm:gap-6">
				{#if progress && progress.time_read >= MINUTES_TO_SEE_STATS * 60}
					<div class="min-w-50 text-center">
						<progress
							value={progress.current_page}
							max={book.pages}
							class="h-2 w-full bg-current/20"
							style:--bar-color|important={progress.current_page < book.pages
								? 'var(--color-status-logo-progress)'
								: 'var(--color-status-logo-done)'}
						>
						</progress>
						<span class="text-sm text-current/80 italic">
							{#if progress.current_page < book.pages}
								~{calculateRemainingTime()} till finished | {progress.pages_read % book.pages}p.
							{:else}
								{#if Math.floor(progress.pages_read / book.pages) > 1}
									{Math.floor(progress.pages_read / book.pages)}x |
								{/if}
								{formatMinToHoursOrMin(progress.time_read / 60)} | {progress.pages_read}p.
							{/if}
							| ~{Math.floor(progress.pages_read / (progress.time_read / 3600))}pph
						</span>
					</div>
				{:else}
					<div class="p-1 text-sm text-nowrap text-current/50 italic">
						(Read {MINUTES_TO_SEE_STATS} min.<br /> to see stats)
					</div>
				{/if}
				<button
					class="h-12 min-w-28 rounded-2xl bg-current/15 font-light shadow-lg
			outline-current/40 hover:outline active:font-normal"
					onclick={() => goto(`${location.pathname}/reading`)}
				>
					<span class="text-2xl text-current/80">Read</span>
				</button>
			</div>
		{/if}
	</span>
	<div class="flex-1">
		<div class="float-left mr-3 min-w-42 sm:min-w-56">
			<CoverWrapper cover={book.cover} title={book.title} />
		</div>
		<div>
			{#if page.params.from === 'shared'}
				<MembersList members={$groupMembers.filter((m) => book.owned_by.includes(m.id))}
				></MembersList>
			{/if}
			<p class="mb-4 text-justify hyphens-auto">
				{#if readMore}
					{book.description}
				{:else}
					{book.description.slice(0, descriptionPreview)}...
				{/if}
				<button class="text-sm text-current/80 italic" onclick={() => (readMore = !readMore)}>
					{#if book.description.length > descriptionPreview}
						{#if readMore}
							See less
						{:else}
							Read more
						{/if}
					{/if}
				</button>
			</p>
		</div>
	</div>

	<footer class="flex justify-center sm:justify-start">
		<div class="text-center text-current/80">
			<p>{book.isbn}</p>
			<p>{book.pages} pages</p>
			<p>{toReadableLanguage(book.language)}</p>
			<p>By {book.publisher} in {book.publishedDate}</p>
		</div>
	</footer>
</div>
