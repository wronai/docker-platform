graph TB
    User[👤 User] --> Browser[🌐 Browser]
    Browser --> Caddy[🔶 Caddy Proxy<br/>Port: 80/443]
    
    %% Keycloak Authentication Flow
    Browser -.->|1. Login Request| Keycloak[🔐 Keycloak<br/>Identity Provider<br/>Port: 8443]
    Keycloak -.->|2. JWT Token| Browser
    
    %% Caddy Routing
    Caddy -->|"/"| Flutter[🎨 Flutter Web<br/>Media Vault UI<br/>Port: 3000]
    Caddy -->|"/api/*"| VaultAPI[🔒 Media Vault API<br/>Go + Fiber + JWT<br/>Port: 8080]
    Caddy -->|"/admin/*"| VaultAdmin[⚙️ Media Vault Admin<br/>Port: 3001]
    Caddy -->|"/auth/*"| Keycloak
    
    %% Authentication Flow
    Flutter -.->|3. API Call + JWT| VaultAPI
    VaultAPI -.->|4. Validate JWT| Keycloak
    VaultAdmin -.->|Admin API + JWT| VaultAPI
    
    %% Media Vault Core
    VaultAPI --> SQLite[(📄 SQLite Database<br/>Vault Metadata)]
    VaultAPI --> VaultStorage[🗄️ Vault Storage<br/>Encrypted Files<br/>Persistent Volume]
    
    %% Processing with Service Account
    VaultAPI -->|Service Account Auth| Analyzer[🔍 Media Vault Analyzer<br/>OpenCV + FFmpeg<br/>Port: 8502]
    Analyzer -.->|Service Token| Keycloak
    Analyzer --> VaultStorage
    
    %% Optional NSFW Detection
    VaultAPI -.->|Optional| NSFW[🤖 NSFW Analyzer<br/>TensorFlow<br/>Port: 8501]
    NSFW -.->|Service Token| Keycloak
    
    %% Keycloak Database
    Keycloak --> KeycloakDB[(🐘 PostgreSQL<br/>Keycloak DB<br/>Port: 5432)]
    
    %% Docker Volumes
    subgraph "💾 Docker Volumes"
        VaultData[Vault Data<br/>./data]
        VaultUploads[Vault Files<br/>./uploads]
        KeycloakData[Keycloak Data<br/>keycloak_db_data]
    end
    
    VaultAPI --> VaultData
    VaultAPI --> VaultUploads
    KeycloakDB --> KeycloakData
    
    %% User Roles and Permissions
    subgraph "👥 Keycloak Realm: media-vault"
        AdminRole[🛡️ vault-admin<br/>Full Access]
        UserRole[👤 vault-user<br/>Basic Access]
        ServiceRole[🤖 vault-analyzer<br/>Service Account]
    end
    
    Keycloak --> AdminRole
    Keycloak --> UserRole
    Keycloak --> ServiceRole
    
    %% Security Flow
    subgraph "🔒 Security Flow"
        direction TB
        Login[1. User Login] --> JWT[2. JWT Token]
        JWT --> APICall[3. API Call with JWT]
        APICall --> Validation[4. JWT Validation]
        Validation --> Access[5. Access Granted/Denied]
    end
    
    %% Styling
    classDef frontend fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef auth fill:#fff3e0,stroke:#f57c00,stroke-width:3px
    classDef vault fill:#4a148c,stroke:#7b1fa2,stroke-width:3px,color:#fff
    classDef storage fill:#fff8e1,stroke:#f57c00,stroke-width:2px
    classDef processing fill:#e8f5e8,stroke:#388e3c,stroke-width:2px
    classDef proxy fill:#ffebee,stroke:#d32f2f,stroke-width:2px
    classDef security fill:#f3e5f5,stroke:#9c27b0,stroke-width:2px
    classDef roles fill:#e0f2f1,stroke:#00695c,stroke-width:2px
    
    class Flutter,VaultAdmin frontend
    class Keycloak,KeycloakDB,Login,JWT,APICall,Validation,Access auth
    class VaultAPI vault
    class SQLite,VaultStorage,VaultData,VaultUploads,KeycloakData storage
    class Analyzer,NSFW processing
    class Caddy proxy
    class AdminRole,UserRole,ServiceRole roles