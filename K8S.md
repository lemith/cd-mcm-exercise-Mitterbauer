# Kubernetes — Concepts & Screenshots

## 1. Scaling
kubectl scale deployment product-catalog-api --replicas=3 -n product-catalog

[Screenshot: kubectl get pods -n product-catalog — 3 Pods Running]

## 2. Health Checks

### Readiness vs Liveness Probe
**Readiness Probe** — ist der Pod bereit Traffic zu empfangen?
- Schlägt fehl → Pod wird aus dem Service entfernt, bekommt keinen Traffic mehr
- Pod läuft weiter, wird aber nicht angesprochen
- Typischer Use-case: App startet noch, DB-Connection wird aufgebaut

**Liveness Probe** — lebt der Pod noch?
- Schlägt fehl → Kubernetes killt den Pod und startet ihn neu
- Typischer Use-case: App ist in Deadlock, reagiert nicht mehr

### Was passiert wenn die Probe fehlschlägt?
- Readiness fehlschlägt → kein Traffic, Pod bleibt aber alive
- Liveness fehlschlägt → Pod wird neugestartet (restart)

### Warum unterschiedliche initialDelaySeconds?
- Liveness braucht mehr Zeit (`initialDelaySeconds` höher) — 
  sonst wird die App neugestartet bevor sie überhaupt fertig 
  gestartet ist
- Readiness kann früher prüfen — nur ob Traffic möglich ist,
  nicht ob die App stabil läuft

## 3. Resource Limits

### Was passiert bei Überschreitung?
- **CPU-Limit überschritten** → Pod wird gedrosselt (throttled), 
  läuft langsamer, wird aber nicht gekillt
- **Memory-Limit überschritten** → Pod wird sofort gekillt (OOMKilled)
  und neugestartet

### Warum requests UND limits?
- **Requests** → Kubernetes reserviert diese Ressourcen für den Pod
  (Scheduling-Grundlage)
- **Limits** → maximale Ressourcen die der Pod verbrauchen darf
- Ohne requests: Kubernetes kann Pods nicht sinnvoll auf Nodes verteilen
- Ohne limits: ein Pod kann alle Ressourcen eines Nodes aufbrauchen
  und andere Pods zum Absturz bringen 