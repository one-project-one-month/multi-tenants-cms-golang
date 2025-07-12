package com.content_management_system.gateway;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.util.List;
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


    public ServiceConnectionLogger(DiscoveryClient discoveryClient) {
        this.discoveryClient = discoveryClient;
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