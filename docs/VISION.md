# ChengetAi Vision Document

## Mission

**Preserve African Knowledge | Democratise Education | Reward Knowledge Creators**

## The Problem We Solve

Across Africa:
- 🎓 Millions of students lack access to quality educational resources
- 📚 Invaluable knowledge created by African educators is lost
- 💡 Teachers create content but have no way to monetise their expertise
- 🌍 African curricula are under-represented in global AI training data
- 🔗 Educational institutions lack infrastructure to preserve knowledge
- 🤖 AI systems learn from Western knowledge, missing African context

## Our Solution

ChengetAi is the digital knowledge infrastructure for Africa - a complete ecosystem where:

### For Students
- Access to curated, AI-powered learning from trusted sources
- Learn in their language and curriculum
- Offline-first for connectivity challenges
- AI tutors grounded in real resources
- Fair, transparent, ad-free learning

### For Teachers
- Own your knowledge permanently
- Get rewarded automatically when students use your content
- Contribute once, earn many times
- Build your professional profile
- No need to upload to multiple platforms

### For Institutions
- Preserve institutional knowledge forever
- Comply with knowledge preservation regulations
- Integrate with ZIMSEC, Cambridge, local curricula
- Gain insights into student learning
- Support teacher professional development

### For Africa
- Preserve African knowledge as a continent
- Represent African curricula in global AI
- Build economic opportunity for African educators
- Reduce dependence on foreign educational platforms
- Strengthen educational infrastructure

## The Seven Layers

### Layer 1: Infrastructure
Robust, scalable, African-centric infrastructure:
- Cloud-agnostic (can run on Liquid Intelligent Technologies, AWS, Azure, local)
- Open source tools (PostgreSQL, Redis, MinIO, NATS)
- Observability built-in
- Security by design

### Layer 2: Platform Services
Shared foundation used by all applications:
- Authentication: JWT + OAuth2
- Users: Identity, profiles, permissions
- Schools: Multi-tenant institutional support
- Wallet: Knowledge economy engine
- Payments: Process revenue distribution
- Notifications: Real-time updates
- Analytics: Learning insights
- Search: Full-text and semantic search
- Feature Flags: Safe experimentation

### Layer 3: Discovery Layer (DSpace)
**The Source of Truth** hosted at repo.dare.co.zw

Knowledge is preserved in DSpace Communities structured by:
```
Africa/
├── Zimbabwe/
│   ├── Primary/
│   ├── Secondary/
│   │   ├── O Level/
│   │   └── A Level/
│   ├── TVET/
│   ├── Polytechnic/
│   └── University/
├── Kenya/
├── South Africa/
└── [Every African Country]
```

Every resource in DSpace becomes:
- ✅ **Permanent** - Preserved forever
- 🔍 **Discoverable** - Searchable across the platform
- 📖 **Citable** - Has persistent identifier (DOI)
- 🎯 **Searchable** - Full-text and semantic search
- 🤖 **AI-Ready** - Processed for intelligence
- 💰 **Monetisable** - Generates revenue for creators

### Layer 4: Knowledge Layer
Intelligent knowledge processing entirely in Go:
- Harvest from DSpace via REST/OAI-PMH
- Extract text from documents (OCR)
- Generate embeddings (vector representations)
- Build semantic indexes
- Create knowledge graphs
- Extract curriculum metadata
- Identify subjects, countries, languages
- Prepare context for AI responses

### Layer 5: Knowledge Economy Layer
**Our biggest innovation** - transforms teachers into entrepreneurs:

```
Resource in DSpace
    ↓
Teacher Contributions
    ↓
Views (visible learning interest)
Downloads (resource usage)
Bookmarks (saved for later)
AI Citations (AI used this in response)
Quiz Usage (testing knowledge)
    ↓
Revenue Generated
    ↓
Automatic Distribution to Contributors
    ↓
Teacher Earnings
```

Teachers earn from:
- Direct usage (students viewing/downloading)
- AI usage (citations in AI responses)
- Quiz usage (tests built on their content)
- Institutional subscriptions (schools paying)

**No middleman, no gatekeepers, no exploitation.**

### Layer 6: Intelligence Layer
AI that learns from Africa, not about Africa:

**Core principle**: Every AI response retrieves knowledge first (RAG)
- Never answer from memory
- Always cite sources from DSpace
- Always link to original resources
- Respect creator attribution

Capabilities:
- **AI Tutor**: Real-time learning support
- **Question Answering**: Deep question understanding
- **Summaries**: Extract key concepts
- **Quiz Generation**: Auto-generate assessments
- **Flashcards**: Spaced repetition learning
- **Essay Feedback**: AI-powered writing support
- **Lesson Planning**: Curriculum-aware planning
- **Voice Tutor**: Audio learning support
- **Recommendations**: Personalized learning paths
- **Translation**: Learn in your language
- **Study Planner**: Organize learning goals

### Layer 7: Experience Layer
Multiple applications, one platform:

**ChengetAi Learn** (First Application)
- AI-powered learning for Zimbabwe students
- Support ZIMSEC, Cambridge curricula
- Support primary, secondary, TVET, polytechnic, university
- Offline-first for connectivity
- Local language support

**Future Applications**
- ChengetAi Teacher - Teacher dashboard
- ChengetAi Library - Knowledge discovery
- ChengetAi Research - Academic collaboration
- ChengetAi Parent - Family learning
- ChengetAi Skills - Professional development
- ChengetAi Admin - Institutional management

## The Knowledge Economy Model

### Value Creation Chain

```
Content Creation
    ↓
Contribution to DSpace
    ↓
AI Processing (Embeddings, Summaries, Graph)
    ↓
User Engagement
    ├── Views
    ├── Downloads
    ├── Bookmarks
    ├── AI Citations
    └── Quiz Usage
    ↓
Revenue Generation
    ├── Student subscriptions
    ├── Institution fees
    └── AI usage
    ↓
Automatic Distribution
    ├── Wallet System
    ├── Smart Contracts (Future)
    └── Transparent Ledger
    ↓
Creator Earnings
```

### Who Earns What

**Teachers**: Earn from their contributed resources
- Every view: $0.001
- Every download: $0.01
- Every AI citation: $0.05
- Every quiz taken: $0.10
- (Rates are configurable by institution)

**Institutions**: Earn from their platform usage
- Subscription fees from students
- Institutional licenses
- Analytics and insights

**Platform**: Earns to sustain operations
- Small percentage of transactions
- Service fees on payments
- Premium institutional features

**Students**: Save money
- Free/low-cost access to quality content
- No tracking, no ads
- Support creators directly

## Why DSpace is the Source of Truth

### The Problem with Databases

Traditional platforms (Moodle, Coursera, etc) store everything in their database:
- When platform shuts down, content dies
- Content is locked into one platform
- No ability to migrate to new technology
- Institutions dependent on platform survival
- Knowledge is not truly preserved

### Why DSpace

DSpace is the world's most trusted digital repository:
- Used by 500+ universities globally
- 30+ years of institutional preservation
- Open source (MIT license)
- OAIS-compliant (Trusted Digital Repository standard)
- Handles versioning and history
- Persistent Identifiers (DOI, Handle)
- OAI-PMH for interoperability
- SWORD for automated ingestion
- Rest assured: your knowledge outlives any company

### ChengetAi's Relationship to DSpace

- **Do NOT duplicate** repository functionality
- **Do NOT store** content in ChengetAi database
- **Do NOT replace** repository capabilities
- **Do CONSUME** knowledge from DSpace
- **Do EXTEND** with AI processing
- **Do ENHANCE** with learning interfaces

```
DSpace
  ↑
  │ (REST API)
  │
ChengetAi Knowledge Layer → Processing → Vector DB → Search Index
  ↑
  │ (Retrieve)
  │
ChengetAi Applications ← Uses knowledge
```

## Roadmap

### Phase 1: Foundation (Current)
- ✅ Architectural design
- [ ] Core platform services
- [ ] DSpace integration
- [ ] Basic knowledge processing
- [ ] ChengetAi Learn MVP

### Phase 2: Intelligence
- [ ] Embedding generation
- [ ] RAG implementation
- [ ] AI Tutor
- [ ] Quiz generation
- [ ] Semantic search

### Phase 3: Economy
- [ ] Wallet system
- [ ] Payment processing
- [ ] Revenue distribution
- [ ] Creator dashboards
- [ ] Analytics

### Phase 4: Scale
- [ ] Multi-country support
- [ ] Multiple curricula
- [ ] Institutional features
- [ ] Advanced AI capabilities

### Phase 5: Impact
- [ ] Millions of learners
- [ ] Thousands of institutions
- [ ] Knowledge preserved
- [ ] Sustainable revenue for creators

## Success Metrics

We measure success by impact:

**Education Impact**
- Students learning per month
- Learning outcome improvements
- Geographic coverage across Africa
- Curriculum coverage

**Economic Impact**
- Revenue generated for teachers
- Number of knowledge creators
- Resources contributed
- Institutions participating

**Knowledge Impact**
- Resources preserved
- Languages supported
- Countries served
- Knowledge graph size

**Platform Impact**
- System uptime (99.99%)
- Query latency (<100ms)
- Search accuracy (99.5%+)
- AI response quality

## Core Values

1. **Knowledge First** - Everything serves knowledge
2. **Creator Respect** - Never exploit creators
3. **Quality Always** - Never compromise on learning outcomes
4. **Transparency** - All systems visible, all rules clear
5. **Scalability** - Built for millions from day one
6. **Preservation** - Knowledge must outlive platforms
7. **Africa First** - Built by and for Africans

## The Long Game

We're not building a startup. We're building **infrastructure for African education** that:
- Lasts 100+ years
- Preserves African knowledge
- Rewards African creators
- Serves African learners
- Remains under African control

Every architectural decision must support this long-term vision.

---

**Chengeta - To Preserve, To Protect, To Keep Safe**
