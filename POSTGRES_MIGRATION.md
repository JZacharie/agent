# Migration PostgreSQL - Guide

Ce document explique comment utiliser PostgreSQL avec l'application stLib.

## Configuration

### Option 1 : Utiliser PostgreSQL (Recommandé avec Docker)

Modifiez le fichier `config.toml` :

```toml
database_type = "postgres"
postgres_host = "localhost"
postgres_port = 5432
postgres_user = "stlib"
postgres_password = "stlib"
postgres_database = "stlib"
postgres_sslmode = "disable"
```

### Option 2 : Utiliser SQLite (Par défaut)

Modifiez le fichier `config.toml` :

```toml
database_type = "sqlite"
```

## Déploiement avec Docker Compose

Le fichier `docker-compose.yml` inclut maintenant un service PostgreSQL :

```bash
cd helpers
docker-compose up -d
```

Cela démarre :
- Un conteneur PostgreSQL avec persistance des données
- L'application stLib configurée pour utiliser PostgreSQL

## Variables d'environnement

Vous pouvez également configurer la base de données via des variables d'environnement :

```bash
export DATABASE_TYPE=postgres
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=stlib
export POSTGRES_PASSWORD=stlib
export POSTGRES_DATABASE=stlib
export POSTGRES_SSLMODE=disable
```

Consultez le fichier `.env.example` pour la liste complète des variables disponibles.

## Vérification

Pour vérifier que PostgreSQL fonctionne correctement :

```bash
# Vérifier les logs de l'application
docker-compose logs -f stlib

# Se connecter à PostgreSQL
docker-compose exec postgres psql -U stlib -d stlib

# Lister les tables
\dt

# Quitter psql
\q
```

## Migration des données

⚠️ **Important** : Cette migration ne transfère pas automatiquement les données existantes de SQLite vers PostgreSQL.

Si vous avez des données dans SQLite que vous souhaitez conserver, vous devrez effectuer une migration manuelle des données.

## Retour à SQLite

Pour revenir à SQLite, modifiez simplement `database_type = "sqlite"` dans `config.toml` ou définissez `DATABASE_TYPE=sqlite` dans les variables d'environnement.
