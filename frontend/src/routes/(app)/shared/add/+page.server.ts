import { PUBLIC_API_URL } from '$env/static/public';
import { checkUser } from '$lib/server/inviteValidity';
import { fail, type Actions } from '@sveltejs/kit';

export const actions = {
	createShared: async (event) => {
		const data = await event.request.formData();
		let name = data.get('name')?.toString() ?? '';
		if (name === '') name = 'Group';
		const userId = event.locals.user?.id ?? '';

		const createForm = new FormData();
		createForm.append('id', userId);
		createForm.append('name', name);

		const createRes = await fetch(`${PUBLIC_API_URL}/groups`, {
			method: 'POST',
			body: createForm
		});

		if (createRes.status != 201) {
			return fail(createRes.status);
		}
		const groupId = await createRes.json();
		const emails = data.getAll('emails[]');
		const inviteForm = new FormData();
		emails.forEach((email) => {
			inviteForm.append('email_to', email);
		});
		inviteForm.append('group_id', groupId);
		inviteForm.append('invited_by', userId);
		const inviteRes = await fetch(`${PUBLIC_API_URL}/groups/invites`, {
			method: 'POST',
			body: inviteForm
		});
		if (inviteRes.status != 200) {
			return fail(inviteRes.status);
		}
		return { success: true, groupId };
	},
	checkUser: async (event) => {
		const email = (await event.request.formData()).get('email') ?? '';
		return checkUser(email, event.locals.user);
	}
} satisfies Actions;
