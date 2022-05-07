import { subDays, subHours, subMinutes } from 'date-fns'
import { MATCH_ANY_PARAMETERS, MockedResponses } from 'wildcard-mock-link'

import { getDocumentNode } from '@sourcegraph/http-client'

import {
    BatchSpecWorkspaceResolutionState,
    WorkspaceResolutionStatusResult,
    BatchSpecImportingChangesetsResult,
    PreviewBatchSpecImportingChangesetFields,
    BatchSpecWorkspacesPreviewResult,
    EditBatchChangeFields,
    PreviewVisibleBatchSpecWorkspaceFields,
    BatchSpecState,
    BatchChangeState,
    BatchSpecWorkspaceState,
    HiddenBatchSpecWorkspaceFields,
    ChangesetSpecType,
    VisibleBatchSpecWorkspaceFields,
    BatchSpecWorkspaceStepFields,
} from '../../../graphql-operations'
import { IMPORTING_CHANGESETS, WORKSPACES, WORKSPACE_RESOLUTION_STATUS } from '../create/backend'

const now = new Date()

export const MOCK_USER_NAMESPACE = {
    __typename: 'User',
    id: 'user1234',
    username: 'my-username',
    displayName: 'My Display Name',
    namespaceName: 'my-username',
    viewerCanAdminister: true,
    url: '/users/my-username',
} as const

export const mockBatchChange = (batchChange?: Partial<EditBatchChangeFields>): EditBatchChangeFields => ({
    __typename: 'BatchChange',
    id: 'testbc1234',
    url: '/batch-changes/my-username/my-batch-change',
    name: 'my-batch-change',
    namespace: MOCK_USER_NAMESPACE,
    description: 'This is my batch change description.',
    viewerCanAdminister: true,
    currentSpec: mockBatchSpec(),
    batchSpecs: {
        nodes: [mockBatchSpec()],
    },
    state: BatchChangeState.OPEN,
    ...batchChange,
})

export const mockBatchSpec = (
    batchSpec?: Partial<EditBatchChangeFields['currentSpec']>
): EditBatchChangeFields['currentSpec'] => ({
    __typename: 'BatchSpec',
    id: '1',
    originalInput: 'name: my-batch-change',
    createdAt: now.toISOString(),
    startedAt: null,
    applyURL: null,
    state: BatchSpecState.PENDING,
    ...batchSpec,
})

export const mockWorkspaceResolutionStatus = (
    status: BatchSpecWorkspaceResolutionState,
    error?: string
): WorkspaceResolutionStatusResult => ({
    node: {
        __typename: 'BatchSpec',
        workspaceResolution: {
            __typename: 'BatchSpecWorkspaceResolution',
            state: status,
            failureMessage: error || null,
        },
    },
})

export const UNSTARTED_RESOLUTION: WorkspaceResolutionStatusResult = {
    node: { __typename: 'BatchSpec', workspaceResolution: null },
}

export const mockPreviewWorkspace = (
    id: number,
    fields?: Partial<PreviewVisibleBatchSpecWorkspaceFields>
): PreviewVisibleBatchSpecWorkspaceFields => ({
    __typename: 'VisibleBatchSpecWorkspace',
    id: `id-${id}`,
    path: '/',
    searchResultPaths: ['/first-path'],
    cachedResultFound: false,
    ignored: false,
    unsupported: false,
    ...fields,
    repository: {
        __typename: 'Repository',
        id: `repo-${id}`,
        name: `github.com/my-org/repo-${id}`,
        url: 'superfake.com',
        ...fields?.repository,
    },
    branch: {
        __typename: 'GitRef',
        id: 'main-branch-id',
        displayName: 'main',
        ...fields?.branch,
        target: {
            __typename: 'GitObject',
            oid: 'asdf1234',
            ...fields?.branch?.target,
        },
        url: 'superfake.com',
    },
})

export const mockPreviewWorkspaces = (count: number): PreviewVisibleBatchSpecWorkspaceFields[] =>
    [...new Array(count).keys()].map(id => mockPreviewWorkspace(id))

export const mockStep = (
    number: number,
    step?: Partial<BatchSpecWorkspaceStepFields>
): BatchSpecWorkspaceStepFields => ({
    __typename: 'BatchSpecWorkspaceStep',
    cachedResultFound: false,
    container: 'ubuntu:18.04',
    diffStat: { __typename: 'DiffStat', added: 10, changed: 5, deleted: 5 },
    environment: [],
    exitCode: 0,
    finishedAt: subMinutes(now, 1).toISOString(),
    ifCondition: null,
    number,
    outputLines: ['stdout: Hello World', 'stdout: '],
    outputVariables: [],
    run: `echo Hello World Step ${number} | tee -a $(find -name README.md)`,
    skipped: false,
    startedAt: subMinutes(now, 2).toISOString(),
    ...step,
})

export const mockWorkspace = (
    numberOfSteps = 1,
    workspace?: Partial<VisibleBatchSpecWorkspaceFields>
): VisibleBatchSpecWorkspaceFields => ({
    __typename: 'VisibleBatchSpecWorkspace',
    id: 'test-1234',
    state: BatchSpecWorkspaceState.COMPLETED,
    searchResultPaths: ['/some/path'],
    queuedAt: subHours(now, 1).toISOString(),
    startedAt: subHours(now, 1).toISOString(),
    finishedAt: now.toISOString(),
    failureMessage: null,
    placeInQueue: null,
    path: '/some/path',
    onlyFetchWorkspace: false,
    ignored: false,
    unsupported: false,
    cachedResultFound: false,
    steps: new Array(numberOfSteps).fill(0).map((_item, index) => mockStep(index + 1)),
    changesetSpecs: [
        {
            description: {
                baseRef: 'main',
                baseRepository: {
                    __typename: 'Repository',
                    name: 'github.com/sourcegraph-testing/batch-changes-test-repo',
                    url: '/github.com/sourcegraph-testing/batch-changes-test-repo',
                },
                body: 'My first batch change!',
                diffStat: { __typename: 'DiffStat', added: 100, changed: 50, deleted: 90 },
                headRef: 'hello-world',
                published: null,
                title: 'Hello World',
                __typename: 'GitBranchChangesetDescription',
            },
            id: 'test-1234',
            type: ChangesetSpecType.BRANCH,
            __typename: 'VisibleChangesetSpec',
        },
    ],
    diffStat: { __typename: 'DiffStat', added: 100, changed: 50, deleted: 90, ...workspace?.diffStat },
    stages: {
        __typename: 'BatchSpecWorkspaceStages',
        setup: [
            {
                command: [],
                durationMilliseconds: 0,
                exitCode: 0,
                key: 'setup.fs',
                out: '',
                startTime: subMinutes(now, 10).toISOString(),
                __typename: 'ExecutionLogEntry',
            },
        ],
        srcExec: {
            command: ['src', 'batch', 'exec', '-f', 'input.json'],
            durationMilliseconds: null,
            exitCode: null,
            key: 'step.src.0',
            out:
                'stdout: {"operation":"PREPARING_DOCKER_IMAGES","timestamp":"2022-04-21T06:26:59.055Z","status":"STARTED","metadata":{}}\nstdout: {"operation":"PREPARING_DOCKER_IMAGES","timestamp":"2022-04-21T06:26:59.055Z","status":"PROGRESS","metadata":{"total":1}}\nstdout: {"operation":"PREPARING_DOCKER_IMAGES","timestamp":"2022-04-21T06:26:59.188Z","status":"PROGRESS","metadata":{"done":1,"total":1}}\nstdout: {"operation":"PREPARING_DOCKER_IMAGES","timestamp":"2022-04-21T06:26:59.188Z","status":"SUCCESS","metadata":{}}\nstdout: {"operation":"DETERMINING_WORKSPACE_TYPE","timestamp":"2022-04-21T06:26:59.188Z","status":"STARTED","metadata":{}}\n',
            startTime: subMinutes(now, 10).toISOString(),
            __typename: 'ExecutionLogEntry',
        },
        teardown: [],
        ...workspace?.stages,
    },
    executor: {
        active: true,
        architecture: 'arm64',
        dockerVersion: '20.10.12',
        executorVersion: '0.0.0+dev',
        firstSeenAt: subDays(now, 10).toISOString(),
        gitVersion: '2.35.1',
        hostname: 'some-super-long-hostname.at-some-address.id-123450123123902304723749827498724',
        id: 'test-1234',
        igniteVersion: '',
        lastSeenAt: subMinutes(now, 2).toISOString(),
        os: 'darwin',
        queueName: 'batches',
        srcCliVersion: '3.38.0',
        __typename: 'Executor',
        ...workspace?.executor,
    },
    ...workspace,
    repository: {
        __typename: 'Repository',
        name: 'github.com/sourcegraph-testing/batch-changes-test-repo',
        url: '/github.com/sourcegraph-testing/batch-changes-test-repo',
        ...workspace?.repository,
    },
    branch: {
        __typename: 'GitRef',
        displayName: 'main',
        ...workspace?.branch,
    },
})

export const QUEUED_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.QUEUED,
    placeInQueue: 2,
    startedAt: null,
    finishedAt: null,
    diffStat: null,
    changesetSpecs: [],
})

export const PROCESSING_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.PROCESSING,
    finishedAt: null,
    diffStat: null,
    changesetSpecs: [],
})

export const SKIPPED_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.SKIPPED,
    queuedAt: null,
    startedAt: null,
    finishedAt: null,
    stages: null,
    executor: null,
    diffStat: null,
    changesetSpecs: null,
    ignored: true,
})

export const UNSUPPORTED_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.SKIPPED,
    queuedAt: null,
    startedAt: null,
    finishedAt: null,
    stages: null,
    executor: null,
    diffStat: null,
    changesetSpecs: null,
    unsupported: true,
})

export const LOTS_OF_STEPS_WORKSPACE = mockWorkspace(20)

export const HIDDEN_WORKSPACE: HiddenBatchSpecWorkspaceFields = {
    __typename: 'HiddenBatchSpecWorkspace',
    id: 'id123',
    queuedAt: subMinutes(now, 10).toISOString(),
    startedAt: subMinutes(now, 8).toISOString(),
    finishedAt: subMinutes(now, 2).toISOString(),
    state: BatchSpecWorkspaceState.COMPLETED,
    diffStat: {
        __typename: 'DiffStat',
        added: 10,
        changed: 2,
        deleted: 5,
    },
    placeInQueue: null,
    onlyFetchWorkspace: false,
    ignored: false,
    unsupported: false,
    cachedResultFound: false,
}

export const FAILED_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.FAILED,
    failureMessage: 'failed to perform src-cli step: command failed',
    steps: [mockStep(1), mockStep(2, { exitCode: 1, diffStat: null })],
})

export const CANCELING_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.CANCELING,
    finishedAt: null,
    diffStat: null,
    changesetSpecs: [],
})

export const CANCELED_WORKSPACE = mockWorkspace(1, {
    state: BatchSpecWorkspaceState.CANCELED,
    finishedAt: null,
    diffStat: null,
    changesetSpecs: [],
})

const mockImportingChangeset = (
    id: number
): PreviewBatchSpecImportingChangesetFields & { __typename: 'VisibleChangesetSpec' } => ({
    __typename: 'VisibleChangesetSpec',
    id: `changeset-${id}`,
    description: {
        __typename: 'ExistingChangesetReference',
        externalID: `external-changeset-${id}`,
        baseRepository: {
            name: `repo-${id}`,
            url: 'superfake.com',
        },
    },
})

export const mockImportingChangesets = (
    count: number
): (PreviewBatchSpecImportingChangesetFields & {
    __typename: 'VisibleChangesetSpec'
})[] => [...new Array(count).keys()].map(id => mockImportingChangeset(id))

export const mockBatchSpecWorkspaces = (workspacesCount: number): BatchSpecWorkspacesPreviewResult => ({
    node: {
        __typename: 'BatchSpec',
        workspaceResolution: {
            __typename: 'BatchSpecWorkspaceResolution',
            workspaces: {
                __typename: 'BatchSpecWorkspaceConnection',
                totalCount: workspacesCount,
                pageInfo: {
                    hasNextPage: workspacesCount > 0,
                    endCursor: 'end-cursor',
                },
                nodes: mockPreviewWorkspaces(workspacesCount),
            },
        },
    },
})

export const mockBatchSpecImportingChangesets = (importsCount: number): BatchSpecImportingChangesetsResult => ({
    node: {
        __typename: 'BatchSpec',
        importingChangesets: {
            __typename: 'ChangesetSpecConnection',
            totalCount: importsCount,
            pageInfo: {
                hasNextPage: importsCount > 0,
                endCursor: 'end-cursor',
            },
            nodes: mockImportingChangesets(importsCount),
        },
    },
})

export const UNSTARTED_CONNECTION_MOCKS: MockedResponses = [
    {
        request: {
            query: getDocumentNode(WORKSPACES),
            variables: MATCH_ANY_PARAMETERS,
        },
        result: { data: mockBatchSpecWorkspaces(0) },
        nMatches: Number.POSITIVE_INFINITY,
    },
    {
        request: {
            query: getDocumentNode(IMPORTING_CHANGESETS),
            variables: MATCH_ANY_PARAMETERS,
        },
        result: { data: mockBatchSpecImportingChangesets(0) },
        nMatches: Number.POSITIVE_INFINITY,
    },
    {
        request: {
            query: getDocumentNode(WORKSPACE_RESOLUTION_STATUS),
            variables: MATCH_ANY_PARAMETERS,
        },
        result: { data: UNSTARTED_RESOLUTION },
        nMatches: Number.POSITIVE_INFINITY,
    },
]

export const UNSTARTED_WITH_CACHE_CONNECTION_MOCKS: MockedResponses = [
    {
        request: {
            query: getDocumentNode(WORKSPACES),
            variables: MATCH_ANY_PARAMETERS,
        },
        result: { data: mockBatchSpecWorkspaces(50) },
        nMatches: Number.POSITIVE_INFINITY,
    },
    {
        request: {
            query: getDocumentNode(IMPORTING_CHANGESETS),
            variables: MATCH_ANY_PARAMETERS,
        },
        result: { data: mockBatchSpecImportingChangesets(20) },
        nMatches: Number.POSITIVE_INFINITY,
    },
    {
        request: {
            query: getDocumentNode(WORKSPACE_RESOLUTION_STATUS),
            variables: MATCH_ANY_PARAMETERS,
        },
        result: { data: mockWorkspaceResolutionStatus(BatchSpecWorkspaceResolutionState.COMPLETED) },
        nMatches: Number.POSITIVE_INFINITY,
    },
]
