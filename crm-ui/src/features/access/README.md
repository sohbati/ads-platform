# Access

Roles, actions, and staff groups for CRM IAM.

Workflows: list, create, delete. Groups are labels only; access still comes from staff roles.

API: `/api/crm/panel/{roles,actions,staff-groups}` → backend `/api/v1/panel/*`.

Permissions: `role.*`, `action.*`, `staff_group.*`.
