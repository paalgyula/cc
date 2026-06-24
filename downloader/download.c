#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <netdb.h>
#include <arpa/inet.h>
#include <unistd.h>

#define BUFFER_SIZE 4096

static void print_usage(const char *prog) {
    fprintf(stderr, "Usage: %s [host] [port] [path] [output]\n", prog);
    fprintf(stderr, "\nDownload a file over HTTP.\n");
    fprintf(stderr, "\nArguments:\n");
    fprintf(stderr, "  host    Remote host (default: google.com)\n");
    fprintf(stderr, "  port    Remote port (default: 80, range: 1-65535)\n");
    fprintf(stderr, "  path    HTTP path (default: /)\n");
    fprintf(stderr, "  output  Output filename (default: auto-derived from path)\n");
    fprintf(stderr, "\nExample:\n");
    fprintf(stderr, "  %s example.com 80 /index.html index.html\n", prog);
}

static const char *basename_from_path(const char *path) {
    const char *p = strrchr(path, '/');
    if (p && *(p + 1)) return p + 1;
    if (p) return "index.html";
    return path;
}

static int find_body(FILE *fp, long *body_start) {
    char buf[BUFFER_SIZE];
    size_t total = 0;
    size_t n;

    rewind(fp);

    while ((n = fread(buf, 1, sizeof(buf), fp)) > 0) {
        for (size_t i = 0; i < n; i++) {
            if (total >= 3 &&
                buf[i] == '\n' &&
                buf[i-1] == '\r' &&
                buf[i-2] == '\n' &&
                buf[i-3] == '\r') {
                *body_start = total + 1;
                return 0;
            }
            total++;
        }
    }

    return -1;
}

int main(int argc, char *argv[]) {
    const char *host = "google.com";
    const char *path = "/";
    int port = 80;
    const char *output = NULL;

    if (argc >= 2 && (strcmp(argv[1], "-h") == 0 || strcmp(argv[1], "--help") == 0)) {
        print_usage(argv[0]);
        return 0;
    }

    if (argc > 1) host = argv[1];
    if (argc > 2) {
        long p = atol(argv[2]);
        if (p < 1 || p > 65535) {
            fprintf(stderr, "error: port must be between 1 and 65535, got %ld\n", p);
            return 1;
        }
        port = (int)p;
    }
    if (argc > 3) path = argv[3];

    if (argc > 4) {
        output = argv[4];
    } else {
        output = basename_from_path(path);
    }

    printf("Downloading http://%s:%d%s -> %s\n", host, port, path, output);

    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0) {
        perror("error: socket");
        return 1;
    }

    struct hostent *server = gethostbyname(host);
    if (server == NULL) {
        fprintf(stderr, "error: failed to resolve %s\n", host);
        close(sockfd);
        return 1;
    }

    struct sockaddr_in server_addr;
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);
    memcpy(&server_addr.sin_addr.s_addr, server->h_addr, server->h_length);

    if (connect(sockfd, (struct sockaddr*)&server_addr, sizeof(server_addr)) < 0) {
        perror("error: connect");
        close(sockfd);
        return 1;
    }

    char request[BUFFER_SIZE];
    int n = snprintf(request, sizeof(request),
                     "GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n",
                     path, host);
    if (n < 0 || (size_t)n >= sizeof(request)) {
        fprintf(stderr, "error: request too long\n");
        close(sockfd);
        return 1;
    }

    if (send(sockfd, request, strlen(request), 0) < 0) {
        perror("error: send");
        close(sockfd);
        return 1;
    }

    FILE *fp = tmpfile();
    if (fp == NULL) {
        perror("error: tmpfile");
        close(sockfd);
        return 1;
    }

    char buffer[BUFFER_SIZE];
    int bytes_received;
    while ((bytes_received = recv(sockfd, buffer, sizeof(buffer), 0)) > 0) {
        fwrite(buffer, 1, bytes_received, fp);
    }

    close(sockfd);

    if (bytes_received < 0) {
        perror("error: recv");
        fclose(fp);
        return 1;
    }

    long body_start;
    if (find_body(fp, &body_start) < 0) {
        fprintf(stderr, "error: no HTTP body found in response\n");
        fclose(fp);
        return 1;
    }

    FILE *out = fopen(output, "wb");
    if (out == NULL) {
        perror("error: fopen output");
        fclose(fp);
        return 1;
    }

    char buf[BUFFER_SIZE];
    fseek(fp, body_start, SEEK_SET);
    size_t r;
    while ((r = fread(buf, 1, sizeof(buf), fp)) > 0) {
        fwrite(buf, 1, r, out);
    }

    fclose(out);
    fclose(fp);

    printf("Done: %s (%ld bytes body)\n", output, ftell(fp) < 0 ? 0L : ftell(fp));
    return 0;
}
