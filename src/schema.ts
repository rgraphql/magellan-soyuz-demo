import { buildSchema } from 'graphql'

// schema is the system schema
export const schema = buildSchema(`
# RootQuery is the root query object.
type RootQuery {
counter: Int
names: [String]
allPeople: [Person]
singlePerson: Person
}

# Person represents an individual.
type Person {
name: String
height: Int
}

schema {
query: RootQuery
}
`)
