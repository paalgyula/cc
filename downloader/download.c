#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <netdb.h>
#include <arpa/inet.h>

#define BUFFER_SIZE 1024

int main() {
    const char* host = "google.com";
    const char* path = "/path/to/file.txt";
    const int port = 80;

    // Create a socket
    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0) {
        perror("Failed to create socket");
        return 1;
    }

    // Get server IP address from hostname
    struct hostent* server = gethostbyname(host);
    if (server == NULL) {
        perror("Failed to get server IP address");
        return 1;
    }

    // Create server address structure
    struct sockaddr_in server_addr;
    bzero((char*)&server_addr, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);
    memcpy(&server_addr.sin_addr.s_addr, server->h_addr, server->h_length);

    // Connect to the server
    if (connect(sockfd, (struct sockaddr*)&server_addr, sizeof(server_addr)) < 0) {
        perror("Failed to connect to server");
        return 1;
    }

    // Send HTTP GET request
    char request[BUFFER_SIZE];
    snprintf(request, BUFFER_SIZE, "GET %s HTTP/1.1\r\nHost: %s\r\n\r\n", path, host);
    if (send(sockfd, request, strlen(request), 0) < 0) {
        perror("Failed to send request");
        return 1;
    }

    // Receive and save the response
    FILE* fp = fopen("file.txt", "wb");
    if (fp == NULL) {
        perror("Failed to create output file");
        return 1;
    }

    char buffer[BUFFER_SIZE];
    int bytes_received;
    while ((bytes_received = recv(sockfd, buffer, BUFFER_SIZE, 0)) > 0) {
        fwrite(buffer, 1, bytes_received, fp);
    }

    fclose(fp);
    close(sockfd);

    return 0;
}
