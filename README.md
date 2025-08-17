Objectif : 

La mise en pratique de l'usage des contextes dans le cadre d'une API REST classique. 

- Créer un serveur HTTP 
- Interruption si Le client annule la requête
- Interruption si Le temps maximum est dépassé
- Context avec transmission de données

Ajout d'un exercice context + Middleware + timeouts : 

# Exercice : API de traitement de commandes

## Objectif
Créer une API en Go combinant **middleware**, **context** et **timeout**.

L’API devra :

1. Ajouter un identifiant de requête (`requestID`) à chaque requête via un middleware.
2. Simuler un traitement long dans un endpoint `/process_order`.
3. Respecter un timeout fixé par `context.WithTimeout`.
4. Retourner le `requestID` dans la réponse pour vérifier le middleware.

---

## Étapes

### 1. Middleware `requestID`
- Ajoute un identifiant unique à chaque requête dans le contexte.
- Exemple :
```go
ctx := context.WithValue(req.Context(), "requestID", reqID)
