# JarCraft

**JarCraft** es una herramienta CLI desarrollada en Go que permite detectar y compilar automáticamente proyectos Java utilizando Maven o Gradle.

Su objetivo es simplificar la compilación de múltiples proyectos Java desde un único ejecutable, permitiendo configurar el comportamiento de cada proyecto mediante archivos YAML externos.

## Características

- Detección automática de proyectos Java.
- Soporte para Maven y Gradle.
- Detección y utilización de Maven Wrapper y Gradle Wrapper.
- Escaneo recursivo de directorios.
- Exclusión configurable de directorios.
- Configuración global y específica por proyecto.
- Selección de versión de Java por proyecto.
- Detección de instalaciones de JDK.
- Configuración de `JAVA_HOME` para cada proceso de compilación sin modificar las variables globales del sistema.
- Posibilidad de omitir los tests.
- Configuración externa mediante YAML.
- Parámetros CLI que permiten sobrescribir la configuración YAML.

## Requisitos

Para utilizar JarCraft necesitas:

- Windows 10/11 (el ejecutable publicado está compilado para Windows x64).
- Java JDK instalado.
- Maven o Gradle instalado, o los wrappers correspondientes dentro de cada proyecto.

> Si un proyecto incluye `mvnw.cmd` o `gradlew.bat`, JarCraft puede utilizarlos sin necesitar una instalación global de Maven o Gradle.

## Instalación

1. Accede a la sección **Releases** del repositorio.
2. Descarga `jarcraft.exe` de la versión que quieras utilizar.
3. Coloca el ejecutable en una ubicación de tu preferencia.
4. Configura tu archivo `jarcraft.yaml`.

Puedes ejecutar JarCraft directamente desde PowerShell.

## Uso

### Consultar la ayuda

```powershell
.\jarcraft.exe --help
```

### Consultar las opciones de compilación

```powershell
.\jarcraft.exe build --help
```

### Compilar proyectos

Para escanear y compilar todos los proyectos Java encontrados dentro de un directorio:

```powershell
.\jarcraft.exe build --path "C:\Users\fern4\Documents\projects"
```

JarCraft detectará los proyectos compatibles y ejecutará el proceso de compilación correspondiente.

### Omitir los tests

```powershell
.\jarcraft.exe build --path "C:\Users\fern4\Documents\projects" --skip-tests
```

### Seleccionar una versión de Java

```powershell
.\jarcraft.exe build --path "C:\Users\fern4\Documents\projects" --java-version 25
```

### Utilizar otro archivo de configuración

```powershell
.\jarcraft.exe build --path "C:\Users\fern4\Documents\projects" --config "C:\Users\fern4\Downloads\jarcraft.yaml"
```

Los parámetros CLI tienen prioridad sobre la configuración del archivo YAML.

## Configuración YAML

JarCraft permite utilizar un archivo externo `jarcraft.yaml` para definir el comportamiento de la herramienta.

Esto permite utilizar el mismo ejecutable en diferentes máquinas con configuraciones distintas.

### Ejemplo de configuración

```yaml
build:
  output: "./builds"
  parallel: false
  fail_fast: false

scanner:
  exclude:
    - ".git"
    - ".idea"
    - ".gradle"
    - "node_modules"

defaults:
  skip_tests: true
  java_version: "25"

projects:
  backend/SpringBootForKafka:
    skip_tests: false
    java_version: "25"

java:
  installations:
    "25": "pathTo25"

logging:
  level: "info"
```

### Opciones de configuración

| Propiedad | Descripción |
|---|---|
| `build.output` | Directorio destinado a los resultados de compilación, si está implementado en el motor. |
| `build.parallel` | Configura la ejecución paralela, si está implementada. |
| `build.fail_fast` | Define si se debe detener el proceso ante un error, si está implementado. |
| `scanner.exclude` | Directorios que se excluirán del escaneo. |
| `defaults.skip_tests` | Define si los tests se omiten globalmente. |
| `defaults.java_version` | Versión de Java predeterminada. |
| `projects` | Configuración específica de cada proyecto. |
| `java.installations` | Rutas de las instalaciones de JDK disponibles. |
| `logging.level` | Nivel de logging configurado. |

> Importante: las opciones que aparecen en el YAML deben estar conectadas a la implementación para producir efectos. La presencia de una propiedad en el archivo no activa automáticamente una funcionalidad.

### Configuración por proyecto

Los proyectos se identifican mediante su ruta relativa al directorio que se está escaneando.

Por ejemplo:

```yaml
defaults:
  skip_tests: true
  java_version: "25"

projects:
  project1:
    skip_tests: false
    java_version: "25"
```

En este ejemplo:

- Todos los proyectos omitirán los tests.
- `SpringBootForKafka` ejecutará sus tests.
- Todos utilizarán Java 25.

### Prioridad de configuración

JarCraft utiliza el siguiente orden de prioridad:

1. Parámetros enviados por CLI.
2. Configuración específica del proyecto.
3. Configuración global (`defaults`).
4. Valores predeterminados internos de JarCraft.

Esto permite adaptar el comportamiento sin modificar el ejecutable.

## Gestión de Java

JarCraft permite registrar diferentes instalaciones de JDK:

```yaml
java:
  installations:
    "17": "pathTo17"
    "21": "pathTo21"
    "25": "pathTo25"
```

Cuando se solicita una versión específica, JarCraft busca una instalación compatible.

Además, configura `JAVA_HOME` y el `PATH` para el proceso de compilación, sin modificar las variables de entorno globales de Windows.

## Compilación de JarCraft desde el código fuente

Si deseas modificar o desarrollar JarCraft, necesitas Go instalado.

Clona el repositorio:

```powershell
git clone <REPOSITORY_URL>
cd jarcraft
```

Descarga las dependencias:

```powershell
go mod download
```

Ejecuta las pruebas:

```powershell
go test -v ./...
```

Compila el ejecutable:

```powershell
go build -o jarcraft.exe .
```

## GitHub Actions y releases

JarCraft utiliza GitHub Actions para automatizar la integración continua y publicación de versiones.

El pipeline contempla:

- Ejecución automática de tests.
- Validación estática mediante `go vet`.
- Compilación para Windows x64.
- Generación de checksum SHA-256.
- Publicación de GitHub Releases.

Los releases se generan cuando se publica un tag que comienza con `v`.

Por ejemplo:

```powershell
git tag v0.1.0
git push origin v0.1.0
```

El pipeline ejecutará las pruebas y, si todo pasa, generará el ejecutable y lo adjuntará al release.

### Versionado semántico

JarCraft utiliza el esquema de versionado semántico:

`MAJOR.MINOR.PATCH`

| Tipo | Ejemplo | Uso |
|---|---|---|
| MAJOR | `v1.0.0` | Cambios incompatibles. |
| MINOR | `v0.2.0` | Nuevas funcionalidades compatibles. |
| PATCH | `v0.1.1` | Correcciones de errores. |

La automatización de versiones mediante Conventional Commits y Release Please es una mejora prevista para el flujo de releases.

## Estructura del proyecto

```text
jarcraft/
├── cmd/
│   ├── root.go
│   └── build.go
├── internal/
│   ├── scanner/
│   ├── builder/
│   ├── engine/
│   ├── config/
│   └── javamanager/
├── test/
│   ├── scanner/
│   ├── builder/
│   ├── engine/
│   ├── config/
│   └── javamanager/
├── .github/
│   └── workflows/
├── jarcraft.yaml
├── main.go
├── go.mod
└── go.sum
```

## Filosofía del proyecto

JarCraft busca ofrecer una forma sencilla, configurable y reproducible de compilar múltiples proyectos Java.

La idea es separar la lógica de compilación de la configuración del entorno, permitiendo que una misma versión del ejecutable se adapte a diferentes máquinas y necesidades.

---

**JarCraft — One CLI. Multiple Java Projects.**