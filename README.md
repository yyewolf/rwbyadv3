# Generic Github Repository Template

Generic github repository template that keeps child repositories sync'd

Use this template as a sensible baseline for new github repositories.

## Instructions

- Create template from repository
- Install the [settings app](https://github.com/apps/settings) on the new repository
- Remove and re-add the `.github/settings.yml` file so the settings app gets enabled
- From the new repository settings page enable "Allow auto-merge"
- Following the [CODEOWNERS SYNTAX](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners#codeowners-syntax) guidelines, update the new repository CODEOWNERS file
- Following our [Github bot guidline documentation](https://konghq.atlassian.net/wiki/spaces/ENGEN/pages/2720268304/How+To+-+Github+Automation+Guidelines) add a github and dependabot secret for AUTO_MERGE_TOKEN
- **Update** the .github/template-sync.yml file in [kong/template-github-release](https://github.com/Kong/template-github-release) repository with the **cloned repository name** to enable template sync changes
- Update .releaserc to have the correct repository name
- Correct the image name in `.github/workflows/release.yaml`
- Correct the image name in `Makefile`
- Remove the sync workflow at `.github/template-sync.yml` and `.github/workflows/sync.yml`
