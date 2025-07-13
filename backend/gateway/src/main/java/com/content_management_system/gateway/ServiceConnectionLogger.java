package com.content_management_system.gateway;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.loadbalancer.LoadBalancerClient;
import org.springframework.context.event.EventListener;
import org.springframework.http.ResponseEntity;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

@Component
public class ServiceConnectionLogger {

    private static final Logger log = LoggerFactory.getLogger(ServiceConnectionLogger.class);
    private final DiscoveryClient discoveryClient;
    private final AtomicBoolean isCmsServiceConnected = new AtomicBoolean(false);
    private static final String CMS_SERVICE_NAME = "cms-service";
    private static final String LMS_SERVICE_NAME = "lms-service";
    private final AtomicBoolean isLmsServiceConnected = new AtomicBoolean(false);
    private static final String CMS_MAIN_SERVICE_NAME = "cms-main-system";
    private final AtomicBoolean isCmsMainServiceConnected = new AtomicBoolean(false);
    private final LoadBalancerClient loadBalancerClient;


    public ServiceConnectionLogger(
        final DiscoveryClient discoveryClient,
        final LoadBalancerClient loadBalancerClient
    ) {
        this.discoveryClient = discoveryClient;
        this.loadBalancerClient = loadBalancerClient;
    }

    @EventListener(ApplicationReadyEvent.class)
    public void onApplicationReady(final ApplicationReadyEvent event) {
        CompletableFuture.delayedExecutor(10, TimeUnit.SECONDS).execute(this::checkServiceConnections);
    }

    @Scheduled(fixedDelay = 30000)
    public void checkServiceConnections() {
        final List<String> services = discoveryClient.getServices();

        for (String serviceName : services) {
            try {
                final ServiceInstance instance = loadBalancerClient.choose(serviceName);
                if (instance != null) {
                    final String uri = instance.getUri().toString();

                    final RestTemplate restTemplate = new RestTemplate();
                    restTemplate.setRequestFactory(new HttpComponentsClientHttpRequestFactory());

                    try {
                        final ResponseEntity<String> response = restTemplate.getForEntity(
                                uri + "/health", String.class);

                        if (response.getStatusCode().is2xxSuccessful()) {
                            log.info("Successfully discovered and connected to '{}' at {}",
                                    serviceName, uri);
                        } else {
                            log.warn("Service '{}' at {} returned status: {}",
                                    serviceName, uri, response.getStatusCode());
                        }
                    } catch (final Exception e) {
                        log.error("Failed to connect to '{}' at {}: {}",
                                serviceName, uri, e.getMessage());
                    }
                } else {
                    log.warn("No instances found for service '{}'", serviceName);
                }
            } catch (final Exception e) {
                log.error("Error checking service '{}': {}", serviceName, e.getMessage());
            }
        }
    }

    @Scheduled(fixedRate = 10000)
    public void checkServiceStatus() {
        checkIndividualService(CMS_SERVICE_NAME, isCmsServiceConnected);
        checkIndividualService(LMS_SERVICE_NAME, isLmsServiceConnected);
        checkIndividualService(CMS_MAIN_SERVICE_NAME, isCmsMainServiceConnected);
    }

    private void checkIndividualService(String serviceName, AtomicBoolean isConnected) {
        try {
            List<ServiceInstance> instances = discoveryClient.getInstances(serviceName);
            if (!instances.isEmpty()) {
                if (isConnected.compareAndSet(false, true)) {
                    ServiceInstance instance = instances.get(0);
                    log.info("Successfully discovered and connected to '{}' at {}:{}",
                            serviceName, instance.getHost(), instance.getPort());
                }
            } else {
                if (isConnected.compareAndSet(true, false)) {
                    log.warn("Connection lost with '{}'. Service is no longer discovered.", serviceName);
                }
            }
        } catch (Exception e) {
            log.error("Error during service discovery check for '{}'", serviceName, e);
            if (isConnected.compareAndSet(true, false)) {
                log.warn("Connection status with '{}' set to disconnected due to an error.", serviceName);
            }
        }
    }
}