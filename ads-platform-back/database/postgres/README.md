# ads platform Database Schema Files

This directory contains the PostgreSQL schema files for the ads platform service, organized by table for better maintainability.

## File Structure

### `0_init.sql`
- Database initialization script
- Creates database and user
- Sets up permissions and schema privileges
- **Must be run as postgres superuser**
 


## Features

- **Modular Design**: Each table in its own file for easy maintenance
- **Proper Dependencies**: Foreign key relationships maintained
- **Performance Optimized**: Indexes on frequently queried columns
- **Automatic Timestamps**: Triggers for updated_at fields
- **Documentation**: Comprehensive comments on tables and columns
- **JSONB Support**: Flexible data storage for facilities and images
- **Geographic Data**: Support for location coordinates
- **Many-to-Many Relationships**: Proper junction table implementation
 