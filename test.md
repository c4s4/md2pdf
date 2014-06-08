% id:       1
% date:     2002-10-25
% title:    Toute toute première fois...
% author:   Michel Casabianca
% email:    michel.casabianca@gmail.com
% keywords: 

# Titre de premier niveau

## Titre de deuxième niveau

Ceci est un paragraphe de test. Ce paragraphe commence par un caractère quelconque, c'est à dire un caractère qui ne fait pas partie de la liste des caractères ayant une signification particulière 1[Texte de la note qui est très long afin de vérifier le bon comportement des notes de bas de page qui doivent rester dans leur boite.].

Deuxième paragraphe, séparé du premier par une ligne vide.

- Premier élément d'une liste à **puces**. Deuxième ligne.
- Deuxième élément d'une *liste*.
- Troisième élément d'une ~~liste~~.
- Quatrième élément d'une `liste`.


1. Liste numérotée. Deuxième ligne.
3. Liste numérotée bis.
4. Liste numérotée ter.


```
Code source. 
N'est pas réarrangé et est en fonte fixe.
```

![test.jpg]()

```
#!/usr/bin/env ruby

class Block

  @@map = {}

  def self.new(lines)
    caracter = lines.first[0..0]
    type = @@map[caracter] || Block
    object = type.allocate()
    object.send(:initialize, lines)
    return object
  end

  def initialize(lines)
    @lines = lines
  end

  def to_html
    "<p>#{@lines.join(" ")}</p>"
  end

end

```

Paragraphe avec du texte **en gras**, *en italique*, ~~souligné~~ et `fonte fixe`. On peut protéger des caractères spéciaux comme * (gras), + (italique), _ (souligné) et ~ (fonte fixe) en les faisant précéder du caractère \. L'antislash lui même peut être écrit par la séquence \\. 

Ceci est le [url_du_lien](texte du lien). Voici une note 2[Texte de la deuxième note.].

Voici maintenant un test de caractères spéciaux (de formatage) HTML : <, >, ', " et & qui doivent apparaître correctement.

Test de ponctuation : ; ? !

