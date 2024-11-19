#let exam(
  title: "Examen",
  course: "Nom du Cours",
  date: none,
  duration: none,
  professor: none,
  school: none,
  logo: none,
  body,
) = {
  // Configuration de base du document
  set document(author: professor, title: title)
  set page(
    margin: (left: 2cm, right: 2cm, top: 2cm, bottom: 2cm),
    numbering: "1/1",
  )
  set text(font: "New Computer Modern", lang: "fr")
  
  // En-tête
  grid(
    columns: (1fr, auto),
    align: (left, right),
    if logo != none {
      image(logo, width: 2cm)
    } else {
      []
    },
    align(right)[
      #text(weight: "bold", school)
    ]
  )
  
  v(1cm)
  
  // Informations de l'examen
  align(center)[
    #block(text(weight: "bold", size: 1.5em, title))
    #v(0.5cm)
    #text(weight: "bold", course)
  ]
  
  v(0.5cm)
  
  // Metadata de l'examen
  grid(
    columns: (1fr, 1fr),
    row-gutter: 0.5em,
    if date != none [*Date:* #date] else [],
    if duration != none [*Durée:* #duration] else [],
    if professor != none [*Professeur:* #professor] else [],
  )
  
  v(1cm)
  
  // Instructions générales
  block(width: 100%, inset: 8pt, radius: 4pt, stroke: 0.5pt + black)[
    *Instructions:*
    - Répondez à toutes les questions dans l'espace prévu à cet effet
    - Écrivez lisiblement et justifiez vos réponses
    - Les documents ne sont pas autorisés sauf mention contraire
  ]
  
  v(1cm)
  
  // Contenu de l'examen
  body
}

// Function helper pour les questions
#let question(number, points, content) = {
  block(width: 100%)[
    #set par(justify: true)
    *Question #number (#points points)*
    #v(0.2cm)
    #content
    #v(0.5cm)
  ]
}

#show: doc => exam(
  title: "Contrôle de Mathématiques",
  course: "Mathématiques - Terminale Spécialité",
  date: "18 novembre 2024",
  duration: "2 heures",
  professor: "Mme Bernard",
  school: "Lycée Victor Hugo",
)[
  {{EXERCISES}}
]