## Benutzung

> Vor dem Benutzen der Anwendung muss die Anwendung entweder mit Docker oder ohne Docker gebaut werden. 
> Siehe spätere Abschnitte.

### Service

Um den Service zu starten, muss die Anwendung mit `tree-service <flags...>` gestartet werden.
Als Flags kann die Adresse und der Port unter dem der Server laufen soll konfiguriert werden.

#### Beispiel

```
tree-service --bind=":8090"
```

Der Service läuft dann unter `localhost:8090`.

### CLI

> Vor der Benutzung des CLI muss der Service gestartet sein.

```
tree-cli <flags...> [action] <arguments...>
```

> [action] muss vorhanden sein!

#### Verfügbare Aktionen [action]
| Aktion | Beschreibung | Beispiel |
| --- | --- | --- |
| `list` | Listet alle vorhandenen Baum IDs auf. | `tree-cli trees` |
| `create-tree` | Erstellt einen neuen Baum, liefert die ID des neuen Baums und ein Sizungstoken, welcher benötigt wird um mit dem `tree-service` zu reden. Außerdem wird die Kapazität jedes Baumblattes übergeben. | `tree-cli create-tree 5` |
| `delete-tree` | Löscht einen Baum. Benötigt die ID des zu löschenden Baums und den dazugehörigen Token. | `tree-cli --id=42 --token=abc123 delete-tree` |
| `insert` | Fügt ein neues Schlüssel-Wert Paar in einen Baum ein. Benötigt Baum ID und Token. | `tree-cli --id=42 --token=abc123 insert 6 "Hallo Welt"` |
| `remove` | Löscht ein bestehendes Schlüssel-Wert Paar in einem Baum. Benötigt Baum ID und Token. | `tree-cli --id=42 --token=abc123 remove 6` |
| `search` | Sucht ein bestehendes Schlüssel-Wert Paar in einem Baum. Benötigt Baum ID und Token. | `tree-cli --id=42 --token=abc123 search 6` |
| `traverse` | Traversiert einen Baum. Benötigt Baum ID und Token. Liefert alle Schlüssel-Wert Paare im Baum in sortierter Reihenfolge. | `tree-cli --id=42 --token=abc123 traverse` |

#### Verfügbare Flags <flag>
| Flag | Beschreibung | Beispiel |
| --- | --- | --- |
| `--bind` | Adresse und Port der CLI. | `--remote=":8091"` CLI läuft unter localhost:8091 |
| `--remote` | Adresse und Port des Services mit dem kommuniziert werden soll. | `--remote=":8090"` Service läuft unter localhost:8090 |
| `--remoteName` | Name des remote Actors des Services (Muss nicht angepasst werden, wenn der Service normal gestartet wird). | `--remote-name="tree-service"` |
| `--id` | ID eines Baumes | `--id=5` |
| `--token` | Token eines Baumes | `--token="abc123"` |
| `--timeout` | Timeout für eine Operation | `--timeout=10s` |

## Ausführen ohne Docker

### Einfach

Für Windows wird die Datei `build.bat` ausgeführt, für Linux die `build.sh`.

Die ausführbaren Dateien `tree-cli.exe` und `tree-service.exe` befindet sich dann im `bin` Ordner.

### Detailliert

#### Go packages installieren

```sh
go get
```

#### Notwendige Go Werkzeuge installieren
```sh
go install github.com/gogo/protobuf/protoc-gen-gogoslick
```

#### Bauen der Messages (Google Protocol Buffers)

Zuerst muss der Protocol Buffer Compiler installiert werden. 
Für Windows uns Linux gibt es Binaries auf [GitHub](https://github.com/protocolbuffers/protobuf/releases/).
Alternativ kann man aber auch einen Paketmanager verwenden (z. B. Chocolatey für Windows mit `choco install protoc`).

Im `messages` Ordner befindet sich eine `build.bat` (Windows) und eine `build.sh` (Linux) Datei, welche man **ausführen** kann um die Message Objekte zu kompilieren.


## Ausführen mit Docker

-   Images bauen

    ```
    make docker
    ```

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    docker run --rm --net actors --name treeservice tree-service \
      --bind="treeservice.actors:8090"
    ```

    Damit das funktioniert, müssen Sie folgendes erst im Tree-Service implementieren:

    -   die `main` verarbeitet Kommandozeilenflags und
    -   der Remote-Actor nutzt den Wert des Flags
    -   wenn Sie einen anderen Port als `8090` benutzen wollen,
        müssen Sie das auch im Dockerfile ändern (`EXPOSE...`)

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    docker run --rm --net actors --name treecli tree-cli --bind="treecli.actors:8091" \
      --remote="treeservice.actors:8090" list
    ```

    Hier sind wieder die beiden Flags `--bind` und `--remote` beliebig gewählt und
    in der Datei `treeservice/main.go` implementiert. `trees` ist ein weiteres
    Kommandozeilenargument, dass z.B. eine Liste aller Tree-Ids anzeigen soll.

    Zum Ausprobieren können Sie den Service dann laufen lassen. Das CLI soll ja jedes
    Mal nur einen Befehl abarbeiten und wird dann neu gestartet.

-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

## Ausführen mit Docker ohne vorher die Docker-Images zu bauen

Nach einem Commit baut der Jenkins, wenn alles durch gelaufen ist, die beiden
Docker-Images. Sie können diese dann mit `docker pull` herunter laden. Schauen Sie für die
genaue Bezeichnung in die Consolenausgabe des Jenkins-Jobs.

Wenn Sie die Imagenamen oben (`treeservice` und `treecli`) durch die Namen aus der
Registry ersetzen, können Sie Ihre Lösung mit den selben Kommandos wie oben beschrieben,
ausprobieren.
