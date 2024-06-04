db.createUser({
    user: 'root',
    pwd: 'example',
    roles: [
        { role: 'dbOwner', db: 'luna' }
    ]
})