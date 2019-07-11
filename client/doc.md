# Mémos

J'initialise un reader sur l'entrée standard qui va capturer le message de l'utilisateur

Je l'enregistre dans une struct Message.Body

Je fais une requete tcp sur mon server go (dans un autre programme) qui intercepte la connexion et met le contenu de ma connexion "conn" dans mon buffer ( je le stocke en ram )

Ensuite j'utilise Fprintln pour réécrire le message envoyé dans le write de mon choix ( la connexion en l'occurence ) ensuite je ferme la connexion
