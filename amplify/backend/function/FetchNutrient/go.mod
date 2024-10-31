module github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/function/FetchNutrient

go 1.18

require (
    github.com/aws/aws-lambda-go v1.41.0
    github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/utils v0.0.0
    github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/services v0.0.0
    github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/machinist v0.0.0
)

replace (
    github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/utils => ../../utils
    github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/services => ../../services
    github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/machinist => ../../machinist
)
