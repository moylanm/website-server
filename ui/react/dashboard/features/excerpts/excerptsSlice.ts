import { createEntityAdapter, createSlice, SerializedError } from '@reduxjs/toolkit';
import { api } from '../api/apiSlice';
import { RootState } from '../../app/store';
import { Excerpt } from './types';

const excerptsAdapter = createEntityAdapter({
	sortComparer: (a: Excerpt, b: Excerpt) => a.id - b.id
});

const initialState = excerptsAdapter.getInitialState({
	status: 'idle',
	statusMessage: '',
	error: null as (SerializedError | null),
	authorField: '',
	workField: '',
	bodyField: ''
});

export const excerptsSlice = createSlice({
	name: 'excerpts',
	initialState,
	reducers: {
		setAuthorField(state, { payload }) {
			state.authorField = payload;
		},
		setWorkField(state, { payload }) {
			state.workField = payload;
		},
		setBodyField(state, { payload }) {
			state.bodyField = payload;
		},
		resetPublishForm(state) {
			state.authorField = '';
			state.workField = '';
			state.bodyField = '';
		},
		resetStatus(state) {
			state.status = 'idle';
			state.statusMessage = '';
			state.error = null;
		}
	},
	extraReducers: (builder) => {
		builder
			.addMatcher(api.endpoints.getExcerpts.matchPending, (state) => {
				state.status = 'loading';
			})
			.addMatcher(api.endpoints.getExcerpts.matchRejected, (state, { error }) => {
				state.status = 'failed';
				state.error = error;
			})
			.addMatcher(api.endpoints.getExcerpts.matchFulfilled, (state, { payload }) => {
				state.status = 'succeeded';
				state.ids = payload.map(excerpt => excerpt.id);
				state.entities = payload.reduce((acc, excerpt) => Object.assign(acc, {[excerpt.id]: excerpt}), {})
			})
			.addMatcher(api.endpoints.publishExcerpt.matchRejected, (state, { error }) => {
				state.status = 'failed';
				state.error = error;
			})
			.addMatcher(api.endpoints.publishExcerpt.matchFulfilled, (state, { payload }) => {
				state.status = 'succeeded';
				state.statusMessage = 'Excerpt successfully created';
				state.ids.unshift(payload.id);
				Object.assign(state.entities, {[payload.id]: payload});
			})
			.addMatcher(api.endpoints.updateExcerpt.matchRejected, (state, { error }) => {
				state.status = 'failed';
				state.error = error;
			})
			.addMatcher(api.endpoints.updateExcerpt.matchFulfilled, (state, { payload }) => {
				state.status = 'succeeded';
				state.statusMessage = 'Excerpt successfully updated';
				Object.assign(state.entities, {[payload.id]: payload});
			})
			.addMatcher(api.endpoints.deleteExcerpt.matchRejected, (state, { error }) => {
				state.status = 'failed';
				state.error = error;
			})
			.addMatcher(api.endpoints.deleteExcerpt.matchFulfilled, (state, { payload }) => {
				state.status = 'succeeded';
				state.statusMessage = 'Excerpt successfully deleted';
				state.ids.splice(state.ids.indexOf(payload), 1);
				delete state.entities[payload];
			})
	}
});

export const {
	setAuthorField,
	setWorkField,
	setBodyField,
	resetPublishForm,
	resetStatus
} = excerptsSlice.actions;

export const {
	selectAll: selectAllExcerpts,
	selectById: selecteExcerptById,
	selectIds: selectExcerptIds
} = excerptsAdapter.getSelectors((state: RootState) => state.excerpts);

export default excerptsSlice.reducer;