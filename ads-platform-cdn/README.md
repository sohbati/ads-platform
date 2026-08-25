# ads-platform-cdn

Static assets + JSON APIs for the ads platform (same layout as ads-platform-ui).

## Structure

```
ads-platform-cdn/
├── cmd/server/                 # Entry point
├── cmd/json_to_csv/            # Utility
├── internal/
│   ├── business/
│   │   ├── category/           # Categories API
│   │   ├── city/               # Cities API
│   │   └── attrschema/         # Attr JSON Schema templates API
│   └── core/
│       ├── config/
│       ├── container/
│       ├── jsonstore/          # Cached JSON file loader
│       └── router/
├── cdn/json/                   # Static JSON files
├── config.example
└── run.sh
```

## Quick start

```bash
./run.sh
# or from repo root: ./3.sh
```

## APIs

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/categories` | Categories from `category.json` |
| GET | `/api/cities` | Cities from `cities.json` |
| GET | `/api/attr-schemas` | Attr JSON Schema templates from `attr-schemas.json` |
| GET | `/json/*` | Raw static JSON files |

Default port: **4000**
