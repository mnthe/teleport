## Documentation Layout

This page is useful for tracking content changes across versions.

Oftentimes, well want to trim or combine content.

It also describes all article subsections and video content in one quick view to make planning easier.

### Blogs and READMEs

- Blogs live in https://github.com/gravitational/web and are hosted live at https://goteleport.com/blog/.
- READMEs live in the root directory of each GitHub repository.

### Meta-Documentation

- Meta Documentation Section
  - [Documentation Quicks Start](https://goteleport.com/docs/docs/)
  - [Documentation Best Practices](https://goteleport.com/docs/docs/best-practices/)

### Products Supported

Our current Product and Reference Documentation topics are divided into thirteen (13) sections supporting the four (4) main products:

1. Teleport Application Access
2. Teleport Server Access
3. Teleport Database Access
4. Teleport Kubernetes Access

Along with a fifth for:

5. Teleport Cloud

### Product and Reference Documentation (Article-Level)

The documentation layout and structure is presently:

- [x] Home (Renamed from documentation)
    - [x] Introduction
    - [x] Adopters
    - [x] Getting Started
    - [x] FAQ
    - [x] Changelog
- [ ] Setup (renamed from infrastructure guides)
    - [x] Installation (moved from Home Section)
    - [x] Docker Compose
    - [x] FAQ (ported from admin manual)
        - [ ] Admin manual subsections split
        - [x] Ports, DNS, ~~PEM~~, SSH
    - [x] Terraform (moved from infrastructure guides)
    - [ ] Helm Guides (moved from Kubernetes)
        - [x] AWS EKS Cluster
        - [x] Google Cloud GKE Cluster
        - [x] Customize Deployment Config
        - [x] Migration From Older Charts
- [ ] Reference Section
    - [x] YAML
        - [ ] Need full config update here
        - [ ] Specific configuration examples
        - [ ] Full YAML example with tables
    - [x] CLI
        - [ ] Need full update here
        - [ ] Common Commands / Cheat Sheet
        - [ ] Teleport
        - [ ] tsh
        - [ ] tctl
    - [ ] Glossary (moved from admin guide)
        - [ ] Link terms in Getting Started Guides to here.
    - [x] Metrics
    - [x] API
        - [x] Client libraries
        - [x] Teleport API Introduction
        - [x] Getting Started
        - [x] API Architecture
- [x] Server access (New section).
  - [x] Introduction (New)
  - [x] Getting Started (New)
  - [x] Guides
    - [x] SSH (moved from user guide)
    - [x] PAM
    - [x] OpenSSH Recording Proxy
    - [x] OpenSSH Client
    - [x] Enhanced Session Recording (moved from Features)
    - [ ] Proxy, NAT, Connecting (think we can factor out a bunch of stuff in the admin guide and move into here)
  - [x] Architecture
    - [ ] Should share some topics with the top-level Architecture/Security article.
    - [x] Bastion Pattern
    - [x] Resource Catalog/Introspection (move all introspection, management items into here)
    - [x] Other common Server Access concepts (greater detail than FAQ)
- [x] Kubernetes Access
    - [x] Introduction
    - [x] Getting Started
    - [x] Guides
        - [x] Multiple Clusters
        - [x] CI/CD
        - [x] Federation
        - [x] Migration
        - [x] Standalone
    - [x] Helm Chart Reference
    - [x] Access Controls
- [x] Database Access Section
    - [x] Introduction
    - [x] Getting Started
    - [x] Guides
        - [x] AWS RDS/Aurora PostgresSQL
        - [x] AWS RDS/Aurora MySQL
        - [x] AWS Redshift PostgresSQL
        - [x] GCP Cloud SQL PostgresSQL
        - [x] Self-Hosted PostgreSQL
        - [x] Self-Hosted MySQL
        - [x] Database GUI Clients
    - [x] Access Controls
    - [x] Architecture
    - [x] Reference
        - [x] Configuration
        - [x] CLI
        - [x] Audit Events
    - [x] FAQ
- [x] Application Access Section
    - [x] Introduction
    - [x] Getting Started
    - [x] Guides
        - [x] Connecting Web Apps
        - [x] Integrating with JWTs
        - [x] Application API Access
    - [x] Access Controls
    - [x] Reference
- [x] Teleport Enterprise Section
    - [x] Introduction
    - [x] Quick Start Guide
    - [x] FedRAMP for SSH & K8s
    - [x] Role-Based Access Control
- [x] Cloud Section
    - [x] Introduction
    - [x] Getting Started
    - [x] Architecture
    - [x] Teleport Cloud FAQ
- [ ] Access Controls Section
    - [x] Introduction
    - [x] Getting Started
    - [x] Guides
        - [x] Role Templates
        - [x] Second Factor - U2F
        - [x] Per-session MFA
        - [x] Dual Authorization
        - [x] Impersonation
        - [ ] Access Requests (New, presumable general guide or landing page to the others)
        - [ ] Identity Managment (New, I propose we move all the certificate management items into here)
        - [ ] SSO Access (New, presumable general guide or landing page to the others - we can maybe put general steps in here and wrap or direct to 3rd party provider specific steps)
    - [ ] Single Sign-On (SSO) - moved from Enterprise
        - [x] Single Sign-On (SSO)
        - [x] Azure Active Directory (AD)
        - [x] Active Directory (ADFS)
        - [x] Google Workspace
        - [x] GitLab
        - [x] OneLogin
        - [x] OIDC
        - [x] Okta
    - [ ] Access Request - moved from Enterprise
        - [x] Integrating Teleport with Slack
        - [x] Integrating Teleport with Mattermost
        - [x] Integrating Teleport with Jira Cloud
        - [x] Integrating Teleport with Jira Server
        - [x] Integrating Teleport with PagerDuty
    - [x] FAQ
- [x] Architecture Section
    - [x] Architecture Overview
    - [x] Teleport Users
    - [x] Teleport Nodes
    - [x] Teleport Auth
    - [x] Teleport Proxy
    - [x] Trusted Clusters
        - [ ] Move all Trusted Cluster content from the Admin Guide into here as well.
    - [ ] Security - new Article
- [x] Preview Section
    - [x] Upcoming Releases