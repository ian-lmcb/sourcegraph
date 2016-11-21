// GENERATED from sourcegraph.schema - DO NOT EDIT

package api

var Schema = `schema {
	query: Query
}

interface Node {
	id: ID!
}

type Query {
	root: Root
	node(id: ID!): Node
}

type Root {
	repository(uri: String!): Repository
	remoteRepositories: [RemoteRepository!]!
	remoteStarredRepositories: [RemoteRepository!]!
}

type Repository implements Node {
	id: ID!
	uri: String!
	description: String!
	commit(rev: String!): CommitState!
	latest: CommitState!
	defaultBranch: String!
	branches: [String!]!
	tags: [String!]!
}

type CommitState {
	commit: Commit
	cloneInProgress: Boolean!
}

type Commit implements Node {
	id: ID!
	sha1: String!
	tree(path: String = "", recursive: Boolean = false): Tree
	file(path: String!): File
	languages: [String!]!
}

type Tree {
	directories: [Directory]!
	files: [File]!
}

type Directory {
	name: String!
	tree: Tree!
}

type File {
	name: String!
	content: String!
}

type RemoteRepository {
	uri: String!
	description: String!
	owner: String!
	name: String!
	httpCloneURL: String!
	language: String!
	fork: Boolean!
	mirror: Boolean!
	private: Boolean!
	createdAt: String!
	pushedAt: String!
	vcsSyncedAt: String!
}`
