import type { components } from './schema';

export type Schemas = components['schemas'];

export type Role = Schemas['Role'];
export type User = Schemas['User'];
export type LoginRequest = Schemas['LoginRequest'];
export type LoginResponse = Schemas['LoginResponse'];
export type ErrorResponse = Schemas['ErrorResponse'];
export type ErrorCode = Schemas['ErrorCode'];

export type Case = Schemas['Case'];
export type CaseStatus = Schemas['CaseStatus'];
export type Complaint = Schemas['Complaint'];
export type Document = Schemas['Document'];
export type CaseDocument = Schemas['Document'];
export type ArchiveEntry = Schemas['ArchiveEntry'];
export type ArchivePatchRequest = Schemas['ArchivePatchRequest'];
export type Verdict = Schemas['Verdict'];

// Trial workflow
export type PreliminaryVerdict = Schemas['PreliminaryVerdict'];
export type CaseOpinion = Schemas['CaseOpinion'];
export type CaseCreateRequest = Schemas['CaseCreateRequest'];
export type CaseAcceptRequest = Schemas['CaseAcceptRequest'];
export type CaseDismissRequest = Schemas['CaseDismissRequest'];
export type OpinionCreateRequest = Schemas['OpinionCreateRequest'];
export type VerdictRequest = Schemas['VerdictRequest'];
export type VerdictResult = Schemas['VerdictResult'];
export type HearingItem = Schemas['HearingItem'];
export type UserListItem = Schemas['UserListItem'];
