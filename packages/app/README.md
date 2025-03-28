# app

This is the main application directory.

Website backend API and application that interfaces with frontend and database. This code is based on guides by Alex Edwards in _Let's Go_ and _Let's Go Further_, as well as `n30w/Darkspace`, with several modifications to fit business requirements and specifications.

## Methods

### Authentication and Authorization

Authentication is done using Auth0 or an external auth provider like Google or GitHub. A persistent user database stores minimal user data, and should only include username and products purchased. Another database is kept during runtime for user session tokens. These tokens define active sessions. Both databases should be GDPR compliant.

## Useful Links

- [What is hexagonal architecture?](https://softengbook.org/articles/hexagonal-architecture)
- [Hexagonal architecture pattern](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/hexagonal-architecture.html)
- [Alistair's Hexagon Pattern](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports and Adapters](https://jmgarridopaz.github.io/content/hexagonalarchitecture.html)
- <https://scalastic.io/en/hexagonal-architecture/>
- <https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/#fundamental-blocks-of-the-system>
