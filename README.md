# fitme-grpc

Rewrite of the FitME rest version to a  gRPC infrastructure with a strong emphasis on middleware, tracing, and logging integration.
The use of a shared proto directory and containerized service architecture gives plenty of flexibility in handling upstream services and dependencies.

The approach to middleware management—ensuring correct ordering for context propagation, handling panics, and managing logs and metrics, wrapping OpenTelemetry's middleware for future-proofing against breaking changes.

TODO
 - Fix Prometheus [x]
 - Loki [x]
 - Tempo [x]
 - Complete Grafana (In production, configure ingres and point to prometheus-grafana)
 - All above with Kubernetes [x]
 - Finish postgres and redis exporter [x]
 - Finish all services
 - Message system to communicate between PT and its clients (add communications between Institution and PTs?)
 - Kafka (?) message queue system and notifications between users and PTs
 - PDF builder: https://github.com/johnfercher/maroto
 - CSV and Excel builder: check Domonda
- K8s deployment
- Where to deploy ?


## Traces
- Exporters: Stdout, Jaeger, Zipkin, Datadog and OpenTelemetry (OTLP) collector
- Importers: OpenTracingShim
## Metrics
- Exporters: Prometheus, Datadog, and OpenTelemetry (OTLP) collector
- Importers: SwiftMetricsShim
## Logs
- Exporters: OpenTelemetry (OTLP) collector

Docker compose just for local testing, telemetry services not relevant in there but all working within Kubernetes.
For dev:
Tempo => kubectl port-forward svc/tempo 4317 -n monitoring
Grafana => kubectl port-forward svc/grafana 3000:80 -n monitoring
Prometheus => kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090 -n monitoring

# Leaderboard feature
=> I am thinking about opening the user platform so every user can see other users progress and plans made etc and introducing a leaderboard system with points for achievements completed. :aPES_Think:
This only for regular users. A "PT" would still only manage their own clients and not see other "PT" clients, obviously.

=> maybe you want progress to be anonymized too. so people can only see the progress in terms of % towards some goal

=> i'm not too sure how embarrassed people are about their fitness goals. but it's one of the main things that keeps them out of gyms from what i can tell reading on the internet

=> :pog: progress leaderboard for friends who train together. can of worms to get into, you want the leaderboards to
have customizable goals (absolute/relative/custom) be public/private  look cute

# FitSynch
FitSynch is a fantastic idea with immense potential in the health and fitness space! By blending personal training management with advanced features like AI-powered assistance for meal planning and shopping, you’re addressing critical pain points for both personal trainers and users. Here's how you could approach and structure the project:

Project Features
Version 1: Core Features

## User Management for Personal Trainers

- Enable trainers to manage their users efficiently.
Features for assigning workout plans, meal plans, and tracking user progress.
Messaging/chat functionality for real-time communication between trainers and users.
Workout Plans

- Allow trainers to create and assign personalized workout plans.
Users can view their daily/weekly workouts with video tutorials or animations for exercises.
Meal Plans

- Create and assign customizable meal plans for users.
Include macros (calories, protein, fats, carbs) for each meal.
Support for different diets (e.g., keto, vegan, gluten-free).
Ingredients and Recipes

- Database of ingredients with nutritional information (calories, protein, etc.).
Recipe creation and sharing between trainers and users.
Automatic generation of shopping lists based on recipes and meal plans.
Shopping List Generator

- Generate a consolidated shopping list from the user's meal plan.
Allow users to check off items as they shop.
Trainer Dashboard

A dedicated dashboard for trainers to manage users, view progress, assign plans, and track payments.

## Version 2: AI-Powered Assistance
- AI for Meal Plans

Suggest meals or recipes based on user preferences, dietary restrictions, and caloric goals.
Dynamically adjust plans based on the user's progress or feedback.
AI for Shopping Assistance

While shopping, users can scan items, and the AI provides:
Calorie breakdown.
Healthier alternatives.
Cost-saving suggestions.
Fitness Insights

AI analyzes user activity, progress, and goals to recommend workout changes or improvements.
Habit Tracking

Track habits like water intake, sleep, or meditation alongside workouts and meals.
AI Chatbot for Instant Help

Provide users with a chatbot to answer fitness and nutrition questions instantly.

**Tech Stack**
_Frontend_
Flutter or React Native: Cross-platform mobile app for iOS and Android.
Web-based dashboard for trainers using React.js or Angular.

_Backend_
Go (Golang): High-performance backend to handle user data, workout plans, and meal plans.
gRPC for communication between services.
PostgreSQL for relational data like user profiles, workouts, and meal plans.
Redis for caching and rate-limiting.
Docker for containerization.

_AI Integration_
Python for AI models (using libraries like TensorFlow or PyTorch).
Integrate AI through REST APIs or gRPC.
Cloud
Use AWS, Google Cloud, or Azure for hosting.
AI models hosted on GPU instances for better performance.

**Key Challenges**
Data Privacy: Ensure the secure handling of sensitive health data.

Use encryption for user data at rest and in transit.
Implement role-based access control (RBAC) for trainers and users.

_AI Training:_

Collect anonymized data for AI training to provide meaningful suggestions.
Partner with nutritionists or fitness experts to train AI models accurately.
Scaling:

Design for scalability to handle a growing user base and complex features like AI.
Use microservices architecture for independent scaling of modules.

## Revenue Model
**Subscription Plans**

For Users: Monthly or annual subscriptions to access premium features like AI assistance and personalized plans.
For Trainers: Tiered pricing based on the number of users they manage.
Marketplace

- Sell premium workout plans, recipes, and AI-generated meal plans.
Affiliate Revenue

- Partnerships with fitness brands, grocery stores, or health products for recommendations.
White-label Offering

- Allow gyms or fitness studios to use the platform under their branding.

**Example User Stories**
User:
Jane wants to lose weight and signs up for FitSynch.
Her trainer assigns her a custom workout and meal plan.
While shopping, Jane scans a product, and the AI suggests a healthier and cheaper alternative.
Jane tracks her progress weekly, and the AI adjusts her goals dynamically.
Trainer:
Mike is a personal trainer managing 15 clients on FitSynch.
He assigns meal plans and workouts through the dashboard, customizing them as needed.
Mike uses AI insights to tweak Jane’s plan after analyzing her weekly progress.

## MVP Launch Plan
**Core Features**: Focus on user management, workout and meal plans, and shopping lists.
**Trainer Dashboard**: Ensure trainers can manage users effectively.
**Scalability**: Prepare infrastructure for AI integrations.
**Marketing**: Launch with gyms and personal trainers as early adopters.

# FitSynch’s Vision
FitSynch combines the best of personal training, nutrition management, and AI. It empowers users to achieve their fitness goals while helping trainers scale their services. With a seamless blend of human expertise and AI-driven insights, it has the potential to disrupt the fitness industry.
