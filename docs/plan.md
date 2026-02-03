# Project Plan (Weeks 7–10) Mini-Moodle

## Team Responsibilities
- **Bauyrzhan**: Student-side functionality, half of Admin logic, backend integration
- **Samat**: Teacher-side functionality, half of Admin logic, documentation and diagrams

Admin responsibilities are intentionally split to ensure equal workload.

---

## Gantt Plan (Weeks 7–10)

Legend: █ = active week

| Task | Owner | W7 | W8 | W9 | W10 | Deliverable |
|---|---|---:|---:|---:|---:|---|
| Repository setup (main/dev/feature branches, access rights) | Both | █ |  |  |  | Git repository ready |
| Base project skeleton (Go monolith, router, health endpoint) | Bauyrzhan | █ |  |  |  | `go run .` works |
| Project proposal & competitor analysis | Samat | █ |  |  |  | `docs/proposal.md` |
| System architecture, Use-Case, ERD, UML diagrams | Samat | █ | █ |  |  | Diagrams in `docs/diagrams/` |
| Database schema & migration scripts | Bauyrzhan |  | █ |  |  | SQL migrations |
| **Student logic** (view courses, enroll, submit assignments, view grades) | Bauyrzhan |  | █ | █ |  | Student REST endpoints |
| **Teacher logic** (create courses, publish assignments, grade submissions) | Samat |  | █ | █ |  | Teacher REST endpoints |
| **Admin logic – Part 1** (user creation, role assignment) | Bauyrzhan |  |  | █ |  | Admin endpoints (users) |
| **Admin logic – Part 2** (course oversight, basic monitoring) | Samat |  |  | █ |  | Admin endpoints (courses) |
| Authentication & role-based access control | Both |  | █ |  |  | Middleware & auth routes |
| Assignment & submission integration (student ↔️ teacher flow) | Both |  |  | █ | █ | End-to-end flow |
| Manual API tests & integration checklist | Samat |  |  |  | █ | Test checklist |
| Build verification & cleanup (no feature polish) | Both |  |  |  | █ | Clean build |

---

## Notes
- Tasks are evenly split between team members.
- Final UI polishing and advanced features are intentionally excluded.
- Focus is on system architecture, module boundaries, and core workflows.
