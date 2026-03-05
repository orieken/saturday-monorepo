# Architecture Diagrams

## Component Overview

```mermaid
graph TB
    subgraph "User Interface"
        UI[Cartridge UI<br/>:9000]
    end
    
    subgraph "Backend Services"
        API[Test Runner Service<br/>:9001]
        PVC[(Reports PVC<br/>1Gi)]
    end
    
    subgraph "Test Environment"
        WEB[Web App<br/>:8000]
        MOCK[Mock API<br/>:8001]
    end
    
    subgraph "Test Execution"
        JOB[Cucumber Job<br/>Pod]
    end
    
    UI -->|POST /api/runs| API
    UI -->|GET /reports/*| API
    API -->|Creates| JOB
    API -->|Mounts| PVC
    JOB -->|Writes to| PVC
    JOB -->|Tests| WEB
    WEB -->|Calls| MOCK
    
    style UI fill:#0ea5e9,stroke:#0284c7,color:#fff
    style API fill:#8b5cf6,stroke:#7c3aed,color:#fff
    style PVC fill:#f59e0b,stroke:#d97706,color:#fff
    style JOB fill:#10b981,stroke:#059669,color:#fff
    style WEB fill:#ec4899,stroke:#db2777,color:#fff
    style MOCK fill:#6366f1,stroke:#4f46e5,color:#fff
```

## Report Flow

```mermaid
sequenceDiagram
    participant User
    participant UI as Cartridge UI
    participant API as Test Runner Service
    participant K8s as Kubernetes
    participant Job as Cucumber Job
    participant PVC as Reports PVC
    
    User->>UI: Click "Run Test"
    UI->>API: POST /api/runs
    API->>K8s: Create Job
    K8s->>Job: Start Pod
    Job->>PVC: Mount /app/reports
    Job->>Job: Run Cucumber Tests
    Job->>PVC: Write HTML Report
    Job->>K8s: Complete
    K8s->>API: Job Status: Succeeded
    API->>API: Update Run Status
    UI->>API: Poll GET /api/runs/{id}
    API-->>UI: Status: passed, reportUrl
    User->>UI: Click "View Report"
    UI->>API: GET /reports/{suite}/{id}/index.html
    API->>PVC: Read Report
    PVC-->>API: HTML Content
    API-->>UI: HTML Report
    UI-->>User: Display Report
```

## Storage Architecture

```mermaid
graph LR
    subgraph "Kubernetes Cluster"
        subgraph "test-runner namespace"
            SVC[test-runner-service<br/>Pod]
            JOB1[Job Pod 1]
            JOB2[Job Pod 2]
            JOB3[Job Pod N]
            
            subgraph "Persistent Storage"
                PVC[reports-pvc<br/>ReadWriteOnce]
                PV[Persistent Volume<br/>1Gi]
            end
        end
    end
    
    SVC -->|Mount /app/reports| PVC
    JOB1 -->|Mount /app/reports| PVC
    JOB2 -->|Mount /app/reports| PVC
    JOB3 -->|Mount /app/reports| PVC
    PVC -->|Bound to| PV
    
    style SVC fill:#8b5cf6,stroke:#7c3aed,color:#fff
    style JOB1 fill:#10b981,stroke:#059669,color:#fff
    style JOB2 fill:#10b981,stroke:#059669,color:#fff
    style JOB3 fill:#10b981,stroke:#059669,color:#fff
    style PVC fill:#f59e0b,stroke:#d97706,color:#fff
    style PV fill:#ef4444,stroke:#dc2626,color:#fff
```

## Directory Structure on PVC

```
/app/reports/
├── final-cucumber-project/
│   ├── run-abc123/
│   │   ├── index.html      ← HTML report
│   │   ├── cucumber.json   ← JSON report
│   │   └── run.log         ← Test logs
│   ├── run-def456/
│   │   ├── index.html
│   │   ├── cucumber.json
│   │   └── run.log
│   └── run-ghi789/
│       ├── index.html
│       ├── cucumber.json
│       └── run.log
└── another-suite/
    └── run-xyz/
        └── ...
```

## Network Flow

```mermaid
graph TB
    subgraph "External Access"
        BROWSER[Browser]
    end
    
    subgraph "Kubernetes Cluster"
        subgraph "NodePort Services"
            NP_UI[NodePort :30000]
            NP_API[NodePort :30001]
            NP_WEB[NodePort :30003]
        end
        
        subgraph "ClusterIP Services"
            SVC_UI[test-runner-ui<br/>:9000]
            SVC_API[test-runner-service<br/>:9001]
            SVC_WEB[web-app<br/>:8000]
            SVC_MOCK[mock-api<br/>:8001]
        end
        
        subgraph "Pods"
            POD_UI[UI Pod]
            POD_API[Service Pod]
            POD_WEB[Web App Pod]
            POD_MOCK[Mock API Pod]
            POD_JOB[Job Pod]
        end
    end
    
    BROWSER -->|http://localhost:30000| NP_UI
    BROWSER -->|http://localhost:30001| NP_API
    NP_UI --> SVC_UI
    NP_API --> SVC_API
    NP_WEB --> SVC_WEB
    SVC_UI --> POD_UI
    SVC_API --> POD_API
    SVC_WEB --> POD_WEB
    SVC_MOCK --> POD_MOCK
    POD_JOB -->|http://web-app:8000| SVC_WEB
    POD_WEB -->|http://mock-api:8001| SVC_MOCK
    
    style BROWSER fill:#0ea5e9,stroke:#0284c7,color:#fff
    style NP_UI fill:#8b5cf6,stroke:#7c3aed,color:#fff
    style NP_API fill:#8b5cf6,stroke:#7c3aed,color:#fff
    style NP_WEB fill:#8b5cf6,stroke:#7c3aed,color:#fff
```

## RBAC Permissions

```mermaid
graph LR
    subgraph "RBAC Configuration"
        SA[ServiceAccount<br/>test-runner-sa]
        ROLE[Role<br/>test-runner-role]
        BINDING[RoleBinding<br/>test-runner-binding]
    end
    
    subgraph "Permissions"
        P1[Create Jobs]
        P2[Get Jobs]
        P3[List Jobs]
        P4[Delete Jobs]
        P5[Get Pods]
        P6[List Pods]
        P7[Get Pod Logs]
    end
    
    subgraph "Pod"
        POD[test-runner-service<br/>Pod]
    end
    
    SA --> BINDING
    ROLE --> BINDING
    BINDING --> POD
    ROLE --> P1
    ROLE --> P2
    ROLE --> P3
    ROLE --> P4
    ROLE --> P5
    ROLE --> P6
    ROLE --> P7
    
    style SA fill:#10b981,stroke:#059669,color:#fff
    style ROLE fill:#f59e0b,stroke:#d97706,color:#fff
    style BINDING fill:#8b5cf6,stroke:#7c3aed,color:#fff
    style POD fill:#ec4899,stroke:#db2777,color:#fff
```

## Deployment Process

```mermaid
graph TB
    START([Start]) --> BUILD[Build Docker Images]
    BUILD --> LOAD[Load Images into Kind]
    LOAD --> NS[Create Namespace]
    NS --> RBAC[Apply RBAC]
    RBAC --> PVC[Create PVC]
    PVC --> SVC[Deploy Service]
    SVC --> UI[Deploy UI]
    UI --> MOCK[Deploy Mock API]
    MOCK --> WEB[Deploy Web App]
    WEB --> NP[Create NodePort Services]
    NP --> WAIT[Wait for Pods Ready]
    WAIT --> DONE([Ready to Test])
    
    style START fill:#10b981,stroke:#059669,color:#fff
    style DONE fill:#10b981,stroke:#059669,color:#fff
    style BUILD fill:#0ea5e9,stroke:#0284c7,color:#fff
    style LOAD fill:#0ea5e9,stroke:#0284c7,color:#fff
    style PVC fill:#f59e0b,stroke:#d97706,color:#fff
```

## Before vs After

### Before (Broken)
```
┌─────────────┐
│ Service Pod │
│             │
│ /app/reports│ ← Empty directory
└─────────────┘

┌─────────────┐
│  Job Pod    │
│             │
│ /reports    │ ← Writes here (emptyDir)
└─────────────┘

❌ Reports not accessible to service!
```

### After (Fixed)
```
┌─────────────┐
│ Service Pod │
│             │
│ /app/reports│ ← Mounts PVC
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ reports-pvc │ ← Shared storage
└─────────────┘
       ▲
       │
┌──────┴──────┐
│  Job Pod    │
│             │
│ /app/reports│ ← Mounts same PVC
└─────────────┘

✅ Reports accessible to both!
```
