package com.content_management_system.gateway;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.context.event.EventListener;
import org.springframework.http.ResponseEntity;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.springframework.util.CollectionUtils;
import org.springframework.util.StringUtils;
import org.springframework.web.client.RestTemplate;

import java.util.List;

@Component
public class ServiceConnectionLogger {

    private static final Logger logger = LoggerFactory.getLogger(ServiceConnectionLogger.class);
    private final DiscoveryClient discoveryClient;
    private final RestTemplate restTemplate;

    private static final String DEFAULT_HEALTH_PATH = "/actuator/health";

    @Autowired
    public ServiceConnectionLogger(final DiscoveryClient discoveryClient) {
        this.discoveryClient = discoveryClient;
        this.restTemplate = new RestTemplate(new HttpComponentsClientHttpRequestFactory());
    }

    @EventListener(ApplicationReadyEvent.class)
    public void onApplicationReady() {
        logger.info("✅ Application is ready. Performing initial health check in 10 seconds...");
        new Thread(() -> {
            try {
                Thread.sleep(10000);
                this.checkServiceConnections();
            } catch (InterruptedException e) {
                logger.error("Initial health check delay was interrupted", e);
                Thread.currentThread().interrupt();
            }
        }).start();
    }

    @Scheduled(fixedDelay = 30000)
    public void checkServiceConnections() {
        logger.info("🔎 Scheduled task running: Checking service connections...");

        final List<String> services = this.discoveryClient.getServices();
        if (CollectionUtils.isEmpty(services)) {
            logger.warn("No services found in Consul discovery. Will retry on the next cycle.");
            return;
        }

        logger.info("Found {} services: {}", services.size(), services);

        services.forEach(serviceId -> {
            final List<ServiceInstance> instances = this.discoveryClient.getInstances(serviceId);
            if (instances.isEmpty()) {
                logger.warn("Service '{}' has no registered instances.", serviceId);
            } else {
                instances.forEach(this::checkInstanceHealth);
            }
        });
    }

    private void checkInstanceHealth(final ServiceInstance instance) {
        String healthPath = instance.getMetadata().get("health-check-path");

        if (!StringUtils.hasText(healthPath)) {
            healthPath = DEFAULT_HEALTH_PATH;
        }

        if (!healthPath.startsWith("/")) {
            healthPath = "/" + healthPath;
        }

        final String healthCheckUrl = instance.getUri() + healthPath;
        final String serviceId = instance.getServiceId();

        try {
            final ResponseEntity<String> response = this.restTemplate.getForEntity(healthCheckUrl, String.class);
            if (response.getStatusCode().is2xxSuccessful()) {
                logger.info("✔️ Successfully connected to '{}' at {}", serviceId, healthCheckUrl);
            } else {
                logger.error("❌ Failed to connect to '{}' at {}. Status: {}, Response: {}",
                        serviceId, healthCheckUrl, response.getStatusCode(), response.getBody());
            }
        } catch (Exception e) {
            logger.error("❌ Error connecting to '{}' at {}: {}", serviceId, healthCheckUrl, e.getMessage());
        }
    }
}