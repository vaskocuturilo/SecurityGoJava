package com.example.java.aop;


import com.example.java.annotation.RateLimiter;
import com.example.java.exception.RateLimitExceededException;
import com.example.java.service.RateLimiterService;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.stereotype.Component;

@Aspect
@Component
@RequiredArgsConstructor
public class RateLimitAspect {

    private final RateLimiterService rateLimiterService;
    private final HttpServletRequest request;

    @Around("@annotation(rateLimit)")
    public Object around(ProceedingJoinPoint joinPoint, RateLimiter rateLimit) throws Throwable {

        final String ip = extractClientIp(request);
        final String methodName = joinPoint.getSignature().toShortString();
        final String customKey = rateLimit.key().isEmpty() ? methodName : rateLimit.key();
        final String finalKey = ip + ":" + customKey;

        boolean allowed = rateLimiterService.isAllowed(
                finalKey,
                rateLimit.limit(),
                rateLimit.duration()
        );

        if (!allowed) {
            throw new RateLimitExceededException(
                    "Rate limit exceeded. Max " + rateLimit.limit() +
                            " requests per " + rateLimit.duration() + "s."
            );
        }

        return joinPoint.proceed();
    }

    private String extractClientIp(HttpServletRequest request) {
        final String forwarded = request.getHeader("X-Forwarded-For");
        if (forwarded != null && !forwarded.isBlank()) {
            return forwarded.split(",")[0].trim();
        }
        return request.getRemoteAddr();
    }
}
