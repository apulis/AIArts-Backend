@Library('apulis-build@master') _

buildPlugin ( {
    repoName = 'AIArtsBackend'
    dockerImages = [
        [
            'compileContainer': '',
            'preBuild':[],
            'imageName': 'apulistech/aiarts-backend',
            'directory': '.',
            'dockerfilePath': 'deployment/Dockerfile',
            'arch': ['amd64','arm64']
        ]
    ]
})
